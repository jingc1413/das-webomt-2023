import {outputPDF} from "./outputA4PDF";
// import domtoimage from 'dom-to-image';
// import { saveAs } from 'file-saver';


function getStyle() {
  let styleContent = `
    #print-container {
      display:block;
      -webkit-print-color-adjust: exact;
      -moz-print-color-adjust: exact;
      color-adjust: exact;
    }
    .print-divider-title {
      height:52px;
    }
    @media print {
      @page {
        margin: 12mm 6mm 12mm 12mm;
        size:A4;
      }
      body {
        margin: 0;
      }
      body > :not(#print-container) {
        display:none;
      }
      html,
      body {
          display: block !important;
      }
      #print-container {
          display: block;
      }
      .el-table .el-table__body-wrapper .el-scrollbar__view {
        display: block !important;
      }
  }`;
  let style = document.createElement("style");
  style.innerHTML = styleContent;
  return style;
}
  

function cleanPrint() {
  let div = document.getElementById('print-container')
  if (div) {
    document.querySelector('body')?.removeChild(div)
  }
}
  

function getContainer(html) {
  cleanPrint()
  let container = document.createElement("div");
  container.setAttribute("id", "print-container");
  // container.innerHTML = html;
  container.appendChild(html)
  return container;
}
  

function getLoadPromise(dom) {
  let imgs = dom.querySelectorAll("img");
  imgs = ([]).slice.call(imgs);

  if (imgs.length === 0) {
    return Promise.resolve();
  }

  let finishedCount = 0;
  return new Promise(resolve => {
    function check() {
      finishedCount++;
      if (finishedCount === imgs.length) {
        resolve();
      }
    }
    imgs.forEach(img => {
      img.addEventListener("load", check);
      img.addEventListener("error", check);
    })
  });
}

/**
 * 
 * @param {*} html 
 * @param {*} title 
 * @param {*} isZoom 
 * @param {*} zoomWidth mm
 * @param {*} zoomHeight mm
 */
export default function printHtml(html, title, isZoom=false, zoomWidth=210, zoomHeight=298) {
    if (!html) {
      return
    }
    // domtoimage.toBlob(html)
    // .then(function (blob) {
    //   saveAs(blob, title+'.png');
    // })
    // .catch(function (error) {
    //     console.error('oops, something went wrong!', error);
    // });
    // return;
    let px1mm = get1mmPx();
    // console.log(px1mm);
    let style = getStyle();
    let container = getContainer(html.cloneNode(true));
    document.body.appendChild(style);
    document.body.appendChild(container);
    if (isZoom) {
      let clientWidth = document.body.offsetWidth;
      let rate = (((zoomWidth - 18)*px1mm)/clientWidth);
      rate = rate >= 1?1:rate;
      let version;
      try {
        version = getBrowserInfo();
      } catch (error) {
        version = 103;
      }
      if (version > 102) {
        container.style.zoom = rate;
      } else {
        container.style.transformOrigin = "top left";
        container.style.transform=`scale(${rate})`;
      }
      let pageHeight = (zoomHeight - 24)*px1mm;
      if (version > 102) {
        pageHeight = pageHeight / rate
      }
      // console.log({rate, clientWidth, pageHeight})
      outputPDF({element: container, rate, pageHeight, version});
      // container.style.backgroundColor = "red";
    }
    const cache = document.title;
    window.addEventListener("beforeprint", function (event) {
      document.title = `${title}`
    });
    window.addEventListener("afterprint", function (event) {
      document.title = `${cache}`
    });
    getLoadPromise(container).then(() => {
      container.style.height = 'auto';
      window.print();
      document.body.removeChild(style);
      document.body.removeChild(container);
    });
}

function get1mmPx() {
  let px = 1;
  let tmpNode = document.createElement('DIV');
  tmpNode.style.cssText = 'width:1mm;height:1mm;visibility:hidden';
  document.body.append(tmpNode);
  let widthStr = getComputedStyle(tmpNode, null).width;
  px = Number(widthStr.substring(0, widthStr.length-2));
  tmpNode.parentNode.removeChild(tmpNode);
  return px;
}

function getBrowserInfo() {
      let versionInfo = navigator.userAgent.toLowerCase() ;
      let version;

        let reg_ie = /msie [\d.]+;/gi, reg_fox = /firefox\/[\d.]+/gi, reg_chrome = /chrome\/[\d.]+/gi , reg_saf = /safari\/[\d.]+/gi ;

        // if(versions.indexOf("msie") > 0)    //IE
        // {
        //     return versions.match(reg_ie) ;
        // }

        // if(versions.indexOf("firefox") > 0)    //firefox
        // {
        //     return versions.match(reg_fox) ;
        // }

        if(versionInfo.indexOf("chrome") > 0)     //Chrome
        {
          version = (versionInfo.match(reg_chrome) + "").replace(/[^0-9.]/ig,"").split('.')[0] ;
        }

        // if(versions.indexOf("safari") > 0 && versions.indexOf("chrome") < 0)  //Safari
        // {
        //     return versions.match(reg_saf) ;
        // }

        return Number(version)
}
