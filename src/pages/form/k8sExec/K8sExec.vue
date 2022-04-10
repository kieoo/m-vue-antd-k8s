<template>
  <div>
    <PageLayout :desc="$t('pageDesc')" ></PageLayout>
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
          <a-collapse v-model="activeKey" v-if="kubeconfigYaml.length>0">
            <a-collapse-panel  key="1" header="kubeconfig 解析结果">
              <p class="text-wrapper" style="font-size: 10px">{{kubeconfigYaml}}</p>
            </a-collapse-panel>
          </a-collapse>
      </a-form-model-item>
      <a-form-model-item label="Clusters">
        <a-select v-model="cluster"
                  mode="multiple"
                  placeholder="please select exec cluster"
                  @change="handleClusterChange">
          <a-select-option v-for="c of this.clusterList" :key="c">
            {{ c.split("/")[1] }}
          </a-select-option>
        </a-select>
      </a-form-model-item>
      <a-form-model-item label="Namespaces/Owners">
        <a-select v-model="nsAndOwner"
                  mode="multiple"
                  placeholder="please select exec namespace and owner"
                  @change="handleNsOwnChange">
          <a-select-option v-for="n in this.selectNsAndOwnerList" :key="n">
            {{ n.split("/")[1] }}/{{ n.split("/")[2] }}/{{n.split("/")[3]}}
          </a-select-option>
        </a-select>
      </a-form-model-item>
      <a-form-model-item label="Containers">
        <a-select v-model="container"
                  mode="multiple"
                  placeholder="please select exec container">
          <a-select-option v-for="i in this.selectContainers" :key="i.con">
            {{ i.con.split("/")[1] }}/{{ i.con.split("/")[2] }}/{{i.con.split("/")[3]}}/{{i.con.split("/")[4]}} <span style='color:indianred'>({{i.pod.length}})</span>
          </a-select-option>
        </a-select>
      </a-form-model-item>
      <a-form-model-item>
        <span slot="label">
          Command&nbsp;
            <a-tooltip title="禁止kill,rm 等危险操作">
              <a-icon type="question-circle-o" />
            </a-tooltip>
        </span>
        <a-input v-model="command" type="textarea" />
        <a-space>
          <a-button type="primary" :loading="running" icon="caret-right" @click="exec">
              Run
          </a-button>
          <a-button type="danger" @click="handleStop" icon="close-circle" v-if="this.running">
            Stop
          </a-button>
        <a-switch v-model="acceptkill" checked-children="killer"  v-show="false"/>
        </a-space>
      </a-form-model-item>
      <a-form-model-item label="Log-Out">
        <a-progress
            :stroke-color="{
                  from: '#108ee9',
                  to: '#87d068',
                   }"
            :percent="running_percent"
            :status="progress_status"
        />
      </a-form-model-item>
      <k-codemirror v-model="log" v-bind:cm-options="cmOption()" class="codemirror"></k-codemirror>
    </a-form-model>
  </a-card>
  </div>
</template>

<script>
import KCodemirror from '@/components/input/Codemirror'
import PageLayout from '@/layouts/PageLayout'
import {KSIABLE} from '@/services/api'
import {request, METHOD} from '@/utils/request'
import Cookie from 'js-cookie'

export default {
  components: {
    KCodemirror,
    PageLayout
  },
  name: 'k8sExec',
  i18n: require('./i18n'),
  data () {
    return {
      page: { desc:"tst",  },
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      form: this.$form.createForm(this),
      formData: new FormData(),
      execData: {},
      uploadFiles: [],
      namespace: "",
      kubeconfigYaml: "",
      kubeconfigYamlBase64: "",
      activeKey: [],
      clusterList: [],
      nsAndOwnerList:[],
      podList: [],
      containers: [],
      selectNsAndOwnerList:[],
      selectContainers:[],
      cluster:[],
      nsAndOwner:[],
      container:[],
      command: "",
      log_dir: "",
      log_info: {},
      log: "",
      showKubeconfig: false,
      acceptkill: false,
      searchLoading: false,
      running: false,
      stop: false,
      running_percent: 0,
      progress_status: "normal"
    }
  },
  methods: {
    upKubeConfRequest(data) {
      this.saveFile(data)
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
    },
    cleanForm: function () {
      this.searchLoading = true
      this.formData.delete('files')
      this.formData.delete('namespace')
      this.kubeconfigYaml = ""
      this.kubeconfigYamlBase64 = ""
      this.cluster = []
      this.nsAndOwner = []
      this.container = []
    },
    async uploadKubeConf() {
      this.cleanForm()

      for ( const up of this.uploadFiles) {
         this.formData.append('files', up.up_file)
      }
      this.formData.append('namespace', this.namespace)
      // console.log(this.formData)
      await request(KSIABLE + "/upKubeconfs",
          METHOD.POST, this.formData).then(res => {
            this.kubeconfigYaml = res.data['kube_config_yaml']
            this.kubeconfigYamlBase64 = res.data['kube_config_yaml_base64']
            Cookie.set("kube_config_yaml_base64", this.kubeconfigYamlBase64, {expires: 1})
            this.podList = res.data['plist']
      }).catch( error => {
        if (error.response) {
          this.$message.error("Get Pods info Failed - " + error.response.data.msg)
        } else {
          this.$message.error("Get Pods info Failed -")
        }
      })
      this.searchLoading = false
      this.updateInfo(this.podList)
    },

    updateInfo: function (podList) {
      this.clusterList = []
      this.nsAndOwnerList = []
      this.containers = []
      for ( const p of podList) {
        this.clusterList.push(p['context'] + '/' + p['cluster'])
        for ( const podInfo of p['pods']) {
          this.nsAndOwnerList.push( p['context'] + '/' + p['cluster'] + '/' + podInfo['namespace'] + '/' + podInfo['owner'])
          for (const container of podInfo['containers']) {

            this.containers.push({con: p['context'] + '/' + p['cluster'] + '/' + podInfo['namespace'] + '/' + podInfo['owner'] + '/' + container, pod: podInfo['pod_names']})
          }
        }
      }
    },
    handleClusterChange: function (value) {
      //this.cluster = value
      this.selectNsAndOwnerList = []
      this.selectContainers = []
      let tmpNsAndOwner = []
      let tmpContainer = []

      // if (value.length === 0) {
      //   this.nsAndOwner = []
      // }
      for (const cluster of value ) {
        for ( const nsAndOwner of this.nsAndOwnerList) {
          if (nsAndOwner.indexOf(cluster) !== -1) {
            this.selectNsAndOwnerList.push(nsAndOwner)
          }
        }
        // 匹配的已选 ns
        for ( let i=0; i< this.nsAndOwner.length; i++) {
          if (this.nsAndOwner[i].indexOf(cluster) !== -1 ) {
             tmpNsAndOwner.push(this.nsAndOwner[i])
          }
        }

        // 匹配的已选 container
        for ( let i=0; i< this.container.length; i++) {
          if (this.container[i].indexOf(cluster) !== -1 ) {
            tmpContainer.push(this.container[i])
          }
        }
      }

      this.nsAndOwner = tmpNsAndOwner
      this.container = tmpContainer
    },
    handleNsOwnChange: function (value) {
      //this.nsAndOwner = value
      this.selectContainers = []
      let tmpContainer = []
      // if (value.length === 0) {
      //   this.container = []
      // }

      for (const ns of value ) {
        for ( const container of this.containers) {
          if (container['con'].indexOf(ns) !== -1) {
            this.selectContainers.push(container)
          }
        }
        // 匹配的已选 container
        for ( let i=0; i< this.container.length; i++) {
          // console.log(this.container)
          if (this.container[i].indexOf(ns) !== -1 ) {
            tmpContainer.push(this.container[i])
          }
        }
      }

      this.container = tmpContainer
    },
    async exec() {
      this.execData = {}

      this.execData["kube_conf_bytes_base64"] = this.kubeconfigYamlBase64
      this.execData["command"] = this.command
      this.execData["accept_kill"] = this.acceptkill
      this.execData["containers"] = []

      this.log_dir = ""
      this.log_info = {}
      this.log = ""
      for (let i = 0; i < this.container.length; i++) {
        let containers = {}
        containers['context'] = this.container[i].split('/')[0]
        containers['cluster'] = this.container[i].split('/')[1]
        containers['ns'] = this.container[i].split('/')[2]
        containers['owner'] = this.container[i].split('/')[3]
        containers['pod_name'] = "all"
        containers['container_name'] = this.container[i].split('/')[4]
        this.execData["containers"].push(containers)
      }
      // console.log(this.execData)
      await request(KSIABLE + "/exec",
          METHOD.POST,
          this.execData)
          .then(res => {
            this.log_dir = res.data.exec_log_dir
            // console.log(this.log_dir)
          }).catch(error => {
            if (error.response) {
              this.$message.error("Commit Exec Failed -" + error.response.data.msg)
            } else {
              this.$message.error("Commit Exec Failed")
            }
          })
      // console.log(this.log_dir)
      this.running_percent = 0
      this.running = true
      this.stop = false
      if (this.log_dir.length > 0) {
        await this.loadLoad(this.log_dir)
      }
      this.running = false
      this.progress_status = "normal"

    },
    async loadLoad(dir) {
      let logParam = {"exec_log_dir": dir}
      let loadErr = false
      let cancel = 0
      do {
        await request(KSIABLE + "/loadLog",
            METHOD.POST,
            logParam)
            .then(res => {
              logParam["exec_log_dir"] = res.data["exec_log_dir"]
              logParam["running"] = res.data["running"]
              logParam["finished"] = res.data["finished"]
              logParam["read"] = res.data["read"]
              this.log_info = res.data
              if (this.log_info["content"].length > 0) {
                this.log = this.log + this.log_info["content"] + "\n"
                // console.log(this.log)
              }
            }).catch(error => {
              if (error.response) {
                this.$message.error("Load Logs Failed -" + error.response.data.msg)
                loadErr = true
              } else {
                this.$message.error("Load Logs Failed")
                loadErr = true
              }
              this.running_percent = (100 * (this.log_info['finished'].length + this.log_info['read'].length) / (this.log_info['finished'].length + this.log_info['read'].length + this.log_info['running'].length)).toFixed(2)
              this.progress_status = "active"
            })
        if (this.stop) {
          cancel++
        }
        if (!this.log_info["success"] && !loadErr && cancel<3) {
          await this.sleep(3000)
        }
      } while (!this.log_info["success"] && !loadErr && cancel<3)

      if (!loadErr && !this.stop) {
        this.running_percent = 100
        this.progress_status = "success"
      } else {
        this.progress_status = "exception"
      }
    },
    handleStop: async function () {
      this.stop = true

      let logParam = {exec_log_dir: this.log_dir}
      await request(KSIABLE + "/cancelLog",
          METHOD.POST,
          logParam).catch( error => {
        if (error.response) {
          this.$message.error("Cancel log Failed - " + error.response.data.msg)
        } else {
          this.$message.error("Cancel log Failed -")
        }
      })
    },
    sleep(ms) {
      return new Promise(resolve => setTimeout(resolve, ms))
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
    },
    cmOption: function () {
      return {
        tabSize: 4,
        mode: 'application/x-sh',
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
  font-family: Arial, monospace;
  height: 600px;
}
/deep/ .ant-collapse .ant-collapse-item .ant-collapse-header {
  padding: 1px 36px;
  font-size: 12px;
}
</style>
