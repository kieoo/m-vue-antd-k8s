<template>
  <codemirror ref="myCm"
              :value="code"
              :options="cmOptions"
              @ready="onCmReady"
              @focus="onCmFocus"
              @input="onCmCodeChange">
  </codemirror>
</template>

<script>

import { codemirror } from 'vue-codemirror'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/base16-dark.css'
import 'codemirror/mode/javascript/javascript.js'
import 'codemirror/mode/python/python.js'

export default {
  components: {
    codemirror
  },
  props: { myCode: String },
  name: "KCodemirror",
  data () {
    return {
      code : "",
      cmOptions: {
        tabSize: 4,
        mode: 'text/x-python',
        theme: 'base16-dark',
        lineNumbers: true,
        line: true,
      }
    }
  },
  methods: {
    onCmReady(cm) {
      console.log('the editor is readied!', cm)
    },
    onCmFocus(cm) {
      console.log('the editor is focus!', cm)
    },
    onCmCodeChange(newCode) {
      console.log('this is new code', newCode)
      // 变更code值, 向上传事件codeByValue
      this.$emit('codeByValue', newCode)
    }
  },
  watch: {
    // 监听pro.myCode 父级的变化, 然后赋值
    myCode: {
      handler (val) {
        console.log('watch father mycode', val)
        this.code = val
      },
      deep: true // 需要监听的数据的深度
      // immediate: true,
    }
  },
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror
    }
  },
  mounted() {
    console.log('this is current codemirror object', this.codemirror)
    // you can use this.codemirror to do something...
  }
}
</script>

<style scoped>

</style>
