<template>
  <codemirror
              ref="myCm"
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
import 'codemirror/theme/lesser-dark.css'
// import 'codemirror/theme/monokai.css'

import 'codemirror/mode/yaml/yaml.js'
import 'codemirror/mode/shell/shell.js'
import 'codemirror/mode/javascript/javascript.js'

//代码折叠文件
require('codemirror/addon/fold/foldcode.js')
require('codemirror/addon/fold/foldgutter.js')
require('codemirror/addon/edit/matchbrackets.js')
require('codemirror/addon/fold/brace-fold.js')
//选中行高亮文件
require('codemirror/addon/selection/active-line.js')
//缩进文件
require('codemirror/addon/fold/indent-fold.js')
//代码只能提示
require('codemirror/addon/hint/show-hint.js')
require('codemirror/addon/hint/anyword-hint.js')
//addon文件夹放的是Code Mirror的功能插件
require('codemirror/addon/fold/comment-fold.js')

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
  props: {
    myCode: String,
    cmOptions: {
      type:Object,
      default: function () {
        return {
          tabSize: 4,
          mode: 'text/x-yaml',
          theme: 'base16-dark',
          lineNumbers: true,
          line: true,
          extraKeys: { // 触发按键
            'Ctrl': 'autocomplete'
          },
          hintOptions: {tables: {
              users: ["name", "score", "birthDate"],
              countries: ["name", "population", "size"]
            }},
          highlightSelectionMatches: { showToken: /\w/, annotateScrollbar: true },
          smartIndent: true, // 自动缩进
          lineWrapping: true, //代码折叠
        }
      }
    }
  },
  name: "KCodemirror",
  data () {
    return {}
  },
  methods: {
    onCmReady() {
      console.log('the editor is readied!')
    },
    onCmFocus() {
      console.log('the editor is focus!')
    },
    onCmCodeChange(newCode) {
      // console.log('this is new code', newCode)
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
    console.log('this is current codemirror object')
    // you can use this.codemirror to do something...
  }
}
</script>

<style scoped>

</style>
