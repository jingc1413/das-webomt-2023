<script>
import { h } from 'vue'
export default {
  name: 'MenuItem',
  functional: true,
  props: {
    icon: {
      type: String,
      default: ''
    },
    title: {
      type: String,
      default: ''
    }
  },
  render() {
    const vnodes = []
    if (this.icon) {
      if (this.icon.includes('i-ep')) {
        vnodes.push(h('i', {class: `${this.icon} sub-el-icon`}))
      } else {
        vnodes.push(h('svg-icon', {'icon-class': `${this.icon}`}))
      }
    }
    if (this.title) {
      vnodes.push(h('span', {'v-slot':'title'}, this.title))
    }
    return vnodes
  }
}
</script>

<style scoped>
.sub-el-icon {
  color: currentColor;
  width: 1em;
  height: 1em;
}
</style>
