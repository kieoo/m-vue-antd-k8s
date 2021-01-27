<template>
  <codemirror ref="myCm"
              :value="myCode"
              :options="cmOptions"
              @ready="onCmReady"
              @focus="onCmFocus"
              @input="onCmCodeChange"
              >
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
  //这里的父类的v-mode=xxx 的值将会传入这个名为myCode 的 prop, 父类也默认监控了事件input。
  // 同时当触发一个input事件并附带一个新的值的时候，这个父类的xxx 的 property 将会被更新
  model: {
    prop: 'myCode',
    event: 'input',
  },
  props: { myCode: String },
  name: "KCodemirror",
  data () {
    return {
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
      // 变更code值, 向上传event.input事件codeByValue
      this.$emit('input', newCode)
    }
  },
  // watch: {
  //   // 监听pro.myCode 父级的变化, 然后赋值
  //   myCode: {
  //     handler (val) {
  //       console.log('watch father mycode', val)
  //       this.code = val
  //     },
  //     deep: true // 需要监听的数据的深度
  //     // immediate: true,
  //   }
  // },
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
