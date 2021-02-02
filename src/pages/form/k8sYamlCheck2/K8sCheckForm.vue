<template>
  <a-card :body-style="{padding: '24px 32px'}" :bordered="false">
    <a-form>
      <a-row class="form-row" :gutter="5">
        <a-col :lg="12" :md="12" :sm="24">
        </a-col>
        <a-col :xl="{span: 12}" :lg="{span: 12}" :md="{span: 12}" :sm="24">
          <a-collapse  accordion :v-show="checkKey==null || checkKey.length>0">
            <a-collapse-panel v-for="(item, index) in checkKey" :key="index"
                              :header="item.check_name" style="background-color:indianred">
              <div class="text-wrapper">{{ item.hints }}</div>
            </a-collapse-panel>
          </a-collapse>
        </a-col>
      </a-row>
      <a-form-item>
        <a-button type="primary" block @click="changeMyYaml()">{{$t('submit')}}</a-button>
      </a-form-item>
    </a-form>
  </a-card>
</template>

<script>
import {request, METHOD} from '@/utils/request'

export default {
  name: 'k8sCheck2',
  i18n: require('./i18n'),
  data () {
    return {
      value: 1,
      myOriCode: "",
      k8sCheckCode: "",
      code1: "",
      checkKey: []
    }
  },
  methods: {
    changeMyYaml: function () {
      //this.code2 = this.code1
      console.log('father is readied!', this.code1)
      request("http://" + location.host.split(":")[0] + ":7002/kcheck",
          METHOD.POST,
          {'ori_yaml': this.code1, 'rule_config':'my_rules.yaml', 'rule_name': 'normal'})
          .then(res => (this.checkKey = res.data))
    },
    getSonCode: function (childV) {
      this.code1 = childV
    },
    changeActiveKey(key) {
      console.log(key)
    }
  },
  computed: {
    desc() {
      return this.$t('pageDesc')
    }
  }
}
</script>

<style lang="less" scoped> //使用less , 代替css
.mirrorl{
  width:50%;
  float:left;
  border:1px solid darkslateblue;
  padding:10px;
}
//pre-wrap值的意思是保留空白并且正常换行。
//white-space各属性值详见这里。
// 其实设置为pre即可使换行符发挥作用，但这时文本在div宽度不足时不会自动换行，
// 而是撞破边界延伸到div外部去，所以还得加上wrap
.text-wrapper {
  white-space: pre-wrap;
}

.codemirror /deep/ .CodeMirror {
  height: 500px;
}
</style>
