<template>
  <a-card :body-style="{padding: '24px 32px'}" :bordered="false">
    <a-form>
      <a-form-item
        :label="$t('title')"
        :labelCol="{span: 7}"
        :wrapperCol="{span: 10}"
      >
        <a-input :placeholder="$t('titleInput')" />
        <k-codemirror></k-codemirror>
      </a-form-item>
      <a-form-item style="margin-top: 24px" :wrapperCol="{span: 10, offset: 7}">
        <a-button type="primary">{{$t('submit')}}</a-button>
        <a-button style="margin-left: 8px">{{$t('save')}}</a-button>
      </a-form-item>
    </a-form>
  </a-card>
</template>

<script>
import KCodemirror from '@/components/input/Codemirror'

export default {
  components: {
    KCodemirror
  },
  name: 'k8sCheck',
  i18n: require('./i18n'),
  data () {
    return {
      value: 1
    }
  },
  computed: {
    desc() {
      return this.$t('pageDesc')
    }
  }
}
</script>

<style scoped>

</style>
