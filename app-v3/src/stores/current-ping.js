import { defineStore } from 'pinia'
import { ElMessage } from 'element-plus';
import apix from '@/api'
import { translator as t } from '@/i18n';


export const useCurrentPing = defineStore('currentPing', {
  state: () => ({
    token: null,
    jobStats: {
      IsRunning: false,
      ws: null,
      messageList : [],
      tempMessageStatus: {
        seq: 0,
        recv: false,
      }
    }
  }),
  getters: {
    getJobStats: (state) => state.jobStats,
    getWsMessageList: (state)=>state.jobStats.messageList
  },
  actions: {
    createPingJob: function (wsInfo) {
      let jobInfo = {...wsInfo};
      if (typeof jobInfo.Address == 'string') {
        jobInfo.Address = [jobInfo.Address];
      }
      return new Promise((resolve) => {
        apix.createPingJob(jobInfo).then((res)=>{
          let {Token, IsRunning} = res;
          if(Token){
            this.token = Token;
            this.jobStats.IsRunning = IsRunning;
            resolve(true);
          }else{
            resolve(false);
            if (res && res.msg) {
              ElMessage.error(res.msg);
            }
          }
        }).catch((e)=>{
          if (e && e.msg) {
            ElMessage.error(e.msg);
          }
          resolve(false);
          console.log('createPingJob error',e);
        })
      })
    },
    initWebsocketJob: function () {
      if (this.jobStats.ws) {
        this.jobStats.ws.close();
      }
      this.jobStats.ws = null;
      this.jobStats.IsRunning = false;
      this.jobStats.messageList = [];
      this.jobStats.tempMessageStatus.recv = false;
      this.jobStats.tempMessageStatus.seq = 0;
    },
    openWebsocket: function (token, {openWebsocketCallback}) {
      this.initWebsocketJob();
      let ws = apix.connectPingJob(token, {
        onMessage: (e) => {
          const msgBody = JSON.parse(e.data);
          let {Type, Time, Data} = msgBody;
          let msg;
          switch (Type) {
            case "PingOnSetup":
              this.jobStats.tempMessageStatus.recv = false;
              break;
            case "PingOnSend": {
              let {Addr, Nbytes, Seq, IPAddr} = Data;
              if (this.jobStats.messageList.length == 0) {
                this.jobStats.messageList.push(`Ping ${Addr}(${IPAddr.IP}) ${Nbytes - 8}(${Nbytes + 20}) bytes of data.`)
              }
              
              let add = Addr;
              if (IPAddr.IP != add) {
                add += `(${IPAddr.IP})`;
              }
              if (this.jobStats.tempMessageStatus.recv == true && this.jobStats.tempMessageStatus.seq != Seq) {
                msg = `From ${add}: icmp_seq=${this.jobStats.tempMessageStatus.seq} Unreachable`;
              } else {
                this.jobStats.tempMessageStatus.recv = true;
              }
              this.jobStats.tempMessageStatus.seq = Seq;
              break;
            }
            case "PingOnRecv": {
              let {Addr, Nbytes, Rtt, Ttl, Seq, IPAddr} = Data;
              let add = Addr;
              if (IPAddr.IP != add) {
                add += `(${IPAddr.IP})`;
              }
              if (this.jobStats.tempMessageStatus.seq == Seq) {
                msg = `${Nbytes} bytes from ${add}: seq=${Seq} ttl=${Ttl} time=${Math.ceil(Rtt/1000000)}ms`;
                this.jobStats.tempMessageStatus.recv = false;
              }
              break;
            }
            case "PingStats": {
              let {Addr, PacketsSent, PacketsRecv, PacketLoss, IPAddr} = Data;
              let {MaxRtt, MinRtt, AvgRtt } = Data;
              let loss = Math.ceil((PacketLoss/PacketsRecv)*10000)/100;

              let add = Addr;
              if (IPAddr.IP != add) {
                add += `(${IPAddr.IP})`;
              }

              if (this.jobStats.tempMessageStatus.recv == true) {
                this.jobStats.messageList.push(`From ${add}: icmp_seq=${this.jobStats.tempMessageStatus.seq} Unreachable`)
                this.jobStats.tempMessageStatus.recv = false;
              }

              msg = `--- ${add} Ping statistics ---
 ${PacketsSent} packets transmitted, ${PacketsRecv} received, ${loss?loss:0}% packet loss,
 rtt min/avg/max = ${Math.ceil(MinRtt/1000000)}/${Math.ceil(AvgRtt/1000000)}/${Math.ceil(MaxRtt/1000000)}`;
              break;
            }
            default:
              break;
          }
          if (msg) {
            this.jobStats.messageList.push(msg);
          }
          if (Type == "PingStats") {
            setTimeout(() => {
              this.jobStats.ws?.close();
              this.jobStats.IsRunning = false;
              this.token = null;
            }, 100);
          }
        },
        onClose: (e) => {
          console.log("websocket onClose", e);
          this.jobStats.IsRunning = false;
        },
        onOpen: (e) => {
          console.log("websocket onOpen", e);
          if (openWebsocketCallback) {
            openWebsocketCallback()
          }
        },
        onError: (e) => {
          console.log("websocket onError", e);
          this.jobStats.IsRunning = false;
        }
      })
      if (ws) {
        this.jobStats.ws = ws;
      }
      return Promise.resolve(true);
    },
    runPingJob: function (token) {
      return new Promise((resolve, reject) => {
        apix.runPingJob(token).then((res)=>{
          this.jobStats.IsRunning = true;
          resolve(true);
        }).catch(()=>{
          reject(false);
        })
      })
    },
    stopPingJob: function (token) {
      return new Promise((resolve, reject) => {
        apix.stopPingJob(token).then((res)=>{
          this.jobStats.ws?.close();
          this.jobStats.IsRunning = false;
          this.token = null;
          resolve(true);
        }).catch(()=>{
          reject(false);
        })
      })
    },

  },
});