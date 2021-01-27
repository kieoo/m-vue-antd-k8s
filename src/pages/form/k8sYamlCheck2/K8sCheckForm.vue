<template>
  <a-card :body-style="{padding: '24px 32px'}" :bordered="false">
    <a-form>
      <a-row class="form-row" :gutter="5">
        <a-col :lg="12" :md="12" :sm="24">
          <k-codemirror v-model="code1"></k-codemirror>
        </a-col>
        <a-col :xl="{span: 12}" :lg="{span: 12}" :md="{span: 12}" :sm="24">
          <a-collapse v-model="activeKey">
            <a-collapse-panel key="1" :header="checkKey">
              <p>{{ text }}</p>
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
import KCodemirror from '@/components/input/Codemirror'
import {request, METHOD} from '@/utils/request'

export default {
  components: {
    KCodemirror,
  },
  name: 'k8sCheck',
  i18n: require('./i18n'),
  data () {
    return {
      value: 1,
      myOriCode: "",
      k8sCheckCode: "",
      code1: "",
      checkKey: "check1",
      text: "",
      activeKey: ['1']
    }
  },
  methods: {
    changeMyYaml: function () {
      //this.code2 = this.code1
      console.log('father is readied!', this.code1)
      request("http://" + location.host.split(":")[0] + ":7002/return",
          METHOD.POST,
          {'code': this.code1}).then(res => (this.text = res.data.toString()))
    },
    getSonCode: function (childV) {
      this.code1 = childV
    }
  },
  computed: {
    desc() {
      return this.$t('pageDesc')
    }
  },
  watch: {
    activeKey(key) {
      console.log(key);
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
</style>
