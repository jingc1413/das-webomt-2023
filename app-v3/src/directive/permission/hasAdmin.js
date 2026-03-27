 /**
 * v-checkAdmin is admin user
 */
 
 import {useAuthStore} from '@/stores/auth'

 export default {
   mounted(el, binding, vnode) {
     const { value } = binding
     
     if (value) {

       const username = useAuthStore().loginUserName;
       
       if (username !== 'admin') {
         el.parentNode && el.parentNode.removeChild(el)
       }
     } else {
       console.log('mounted value', value)
      //  throw new Error(`Set the value of the operation permission tag`)
     }
   }
 }
 