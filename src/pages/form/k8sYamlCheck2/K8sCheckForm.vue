<template>
  <a-card :body-style="{padding: '24px 32px'}" :bordered="false">
    <a-form>
      <a-form-item>
        <a-upload-dragger name="file"
                          :multiple="true"
                          :fileList="downloadFiles"
                          :showUploadList="false"
                          :remove="handleDownloadFileRemove"
                          :customRequest="downloadFilesCustomRequest">
          <p class="ant-upload-drag-icon">
            <a-icon type="copy" />
          </p>
          <p class="ant-upload-text">
            {{$t('uploadTitle')}}
          </p>
          <p class="ant-upload-hint">
            <span v-bind:style="{'white-space':'pre-wrap'}">{{$t('uploadContext')}}</span>
          </p>
        </a-upload-dragger>
        <a-upload  directory
                   name="file"
                   :multiple="true"
                   :fileList="downloadFiles"
                   :remove="handleDownloadFileRemove"
                   :customRequest="downloadFilesCustomRequest">
          <a-button> <a-icon type="inbox" /> Upload Chart  Directory </a-button>
        </a-upload>
        <a-button type="primary" icon="poweroff" :loading="iconLoading" @click="commitHelm">
          Helm Change
        </a-button>
      </a-form-item>
      <a-row class="form-row" :gutter="5">
        <a-col :lg="12" :md="12" :sm="24">
          <k-codemirror v-model="code1" class="codemirror"></k-codemirror>
        </a-col>
        <a-col :xl="{span: 12}" :lg="{span: 12}" :md="{span: 12}" :sm="24">
          <a-collapse  accordion :v-if="checkKey!=null && checkKey.length>0">
            <a-collapse-panel v-for="(item, index) in checkKey" :key="index"
                              :header= item.check_name
                              :style="{'background-color': (item.hints.length>0 ? 'indianred' : 'limegreen')}">
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
      form: this.$form.createForm(this),
      formData: new FormData(),
      downloadFiles: [],
      value: 1,
      code1: "",
      checkBody: {},
      checkKey: [],
      checkResult: "",
      iconLoading: false,
    }
  },
  methods: {
    downloadFilesCustomRequest(data) {
      this.saveFile(data)
      console.log(this.downloadFiles )
    },
    saveFile (data) {
      // push message to download list
      let file = this.fileFormatter(data)
      this.downloadFiles.push(file)

    },
    fileFormatter(data) {
      return {
        uid: Math.ceil(Math.random()*1000),    // 文件唯一标识，建议设置为负数，防止和内部产生的 id 冲突
        name: data.file.name,   // 文件名
        status: 'done', // 状态有：uploading done error removed
        up_file: data.file
      }
    },
    handleDownloadFileRemove(file) {
      const index = this.downloadFiles.indexOf(file)
      const newFileList = this.downloadFiles.slice()
      newFileList.splice(index, 1)
      this.downloadFiles = newFileList

      console.log(this.downloadFiles )
    },
    async commitHelm() {
      this.iconLoading = true
      // let config = {
      //   headers :{
      //     'Content-Type': 'multipart/form-data'
      //   }
      // }
      this.formData.delete('files')
      for ( const up of this.downloadFiles) {
         this.formData.append('files', up.up_file)
      }
      // console.log(this.formData)
      await request("http://" + location.host.split(":")[0] + ":7001/upload",
          METHOD.POST, this.formData).then(res => {
            this.code1 = res.data.chart
      }).catch( error => {
        if (error.response) {
          this.$message.error("Change Failed - " + error.response.data.msg)
        } else {
          this.$message.error("Change Failed")
        }
      })
      console.log(this.code1)
      this.iconLoading = false
    },
    changeMyYaml: function () {
      //this.code2 = this.code1
      console.log('father is readied!', this.code1)
      request("http://" + location.host.split(":")[0] + ":7001/kcheck",
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
