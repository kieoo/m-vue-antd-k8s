<template>
  <a-card :body-style="{padding: '24px 32px'}" :bordered="false">
    <a-form-model :model="form" :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-model-item>
        <span slot="label">
        KubeConfig&nbsp;
          <a-tooltip title="支持上传多份配置;  如果用户没有查询ns权限, 需要指定namespace">
            <a-icon type="question-circle-o" />
          </a-tooltip>
        </span>
        <a-space>
          <a-upload  name="file"
                     :multiple="true"
                     :showUploadList="false"
                     :fileList="uploadFiles"
                     :remove="handleUpFileRemove"
                     :customRequest="upKubeConfRequest">
            <a-button> <a-icon type="inbox" /> Upload KubeConf </a-button>
          </a-upload>
          <a-input-search
              v-model = "namespace"
              placeholder="指定namespace(可选)"
              enter-button="Search"
              size="default"
              :loading="searchLoading"
              @search="uploadKubeConf"
          />
        </a-space>
          <a-collapse v-model="activeKey">
            <a-collapse-panel  key="1" header="kubeconfig 解析结果">
              <p class="text-wrapper" style="font-size: 10px">{{kubeconfigYaml}}</p>
            </a-collapse-panel>
          </a-collapse>
      </a-form-model-item>
      <a-form-model-item label="Clusters">
        <a-select
            mode="multiple"
            placeholder="please select exec cluster"
        >
<!--          <a-select-option v-for="i in formData.get('clusters')" :key="i">-->
<!--            {{ formData.get('clusters')[i] }}-->
<!--          </a-select-option>-->
          <a-select-option v-for="i in podList" :key="i.cluster">
            {{ i.cluster }}
          </a-select-option>
        </a-select>
      </a-form-model-item>
      <a-form-model-item label="Namespaces/Owners">
        <a-select v-model="form.ns_own" placeholder="please select exec namespace and owner">
        </a-select>
      </a-form-model-item>
      <a-form-model-item label="Containers">
        <a-select v-model="form.containers" placeholder="please select exec container">
        </a-select>
      </a-form-model-item>
      <a-form-model-item label="Command">
          <a-input v-model="form.desc" type="textarea" />
          <a-button type="primary" :loading="running" @click="running = true">
            Run - >
          </a-button>
          <a-progress
              :stroke-color="{
                  from: '#108ee9',
                  to: '#87d068',
                   }"
              :percent="99.9"
              status="active"
          />
      </a-form-model-item>
      <a-form-model-item label="Log-Out">
        <k-codemirror v-model="code1" class="codemirror"></k-codemirror>
      </a-form-model-item>
    </a-form-model>
  </a-card>
</template>

<script>
import KCodemirror from '@/components/input/Codemirror'
import {request, METHOD} from '@/utils/request'
import Cookie from 'js-cookie'

export default {
  components: {
    KCodemirror,
  },
  name: 'k8sExec',
  i18n: require('./i18n'),
  data () {
    return {
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      form: this.$form.createForm(this),
      formData: new FormData(),
      uploadFiles: [],
      namespace: "",
      kubeconfigYaml: "",
      kubeconfigYamlBase64: "",
      activeKey: [],
      clusterList: [],
      podList: [],
      value: 1,
      code1: "",
      checkBody: {},
      checkKey: [],
      checkResult: "",
      showKubeconfig: false,
      searchLoading: false,
    }
  },
  methods: {
    upKubeConfRequest(data) {
      this.saveFile(data)
      console.log(this.uploadFiles )
    },
    saveFile (data) {
      // push message to download list
      let file = this.fileFormatter(data)
      this.uploadFiles.push(file)

    },
    fileFormatter(data) {
      return {
        uid: Math.ceil(Math.random()*1000),    // 文件唯一标识，建议设置为负数，防止和内部产生的 id 冲突
        name: data.file.name,   // 文件名
        status: 'done', // 状态有：uploading done error removed
        up_file: data.file
      }
    },
    handleUpFileRemove(file) {
      const index = this.uploadFiles.indexOf(file)
      const newFileList = this.uploadFiles.slice()
      newFileList.splice(index, 1)
      this.uploadFiles = newFileList

      console.log(this.uploadFiles )
    },
    cleanForm: function () {
      this.searchLoading = true
      this.formData.delete('files')
      this.formData.delete('namespace')
      this.kubeconfigYaml = ""
      this.kubeconfigYamlBase64 = ""
      this.formData.delete('')
    },
    async uploadKubeConf() {
      this.cleanForm()

      for ( const up of this.uploadFiles) {
         this.formData.append('files', up.up_file)
      }
      this.formData.append('namespace', this.namespace)
      // console.log(this.formData)
      await request("http://" + location.host.split(":")[0] + ":7001/ksiable/upKubeconfs",
          METHOD.POST, this.formData).then(res => {
            this.kubeconfigYaml = res.data.kube_config_yaml
            this.kubeconfigYamlBase64 = res.data.kube_config_yaml_base64
            Cookie.set("kube_config_yaml_base64", this.kubeconfigYamlBase64, {expires: 1})
            this.podList = res.data.plist
      }).catch( error => {
        if (error.response) {
          this.$message.error("Get Pods info Failed - " + error.response.data.msg)
        } else {
          this.$message.error("Get Pods info Failed -")
        }
      })
      console.log(this.podList)
      this.searchLoading = false
      this.activeKey = ['1']

      this.updateInfo(this.podList)
    },

    updateInfo: function (podList) {
      this.clusterList = []
      for (const podinfo of podList) {
        this.clusterList.append(podinfo.cluster)
      }
    },
    changeMyYaml: function () {
      //this.code2 = this.code1
      console.log('father is readied!', this.code1)
      request("http://" + location.host.split(":")[0] + "/kc/kcheck",
          METHOD.POST,
          {'ori_yaml': this.code1, 'rule_config':'default.yaml', 'rule_name': 'deployment'})
          .then(res => {
            this.checkBody = res.data
            this.checkKey = res.data.hints
          }).catch( error => {
            if (error.response) {
              this.$message.error("Change Failed - " + error.response.data.msg)
            } else {
              this.$message.error("Change Failed")
            }
          })
      console.log('father is readied!', this.checkBody)
    },
    changeActiveKey(key) {
      console.log(key)
    },
    collapseColor(level) {
      if (level === 2) {
        return "indianred"
      }
      if (level === 1) {
        return "khaki"
      }
      if (level === 0) {
        return "limegreen"
      }
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
  height: 600px;
}
/deep/ .ant-collapse .ant-collapse-item .ant-collapse-header {
  padding: 1px 36px;
  font-size: 12px;
}
</style>
