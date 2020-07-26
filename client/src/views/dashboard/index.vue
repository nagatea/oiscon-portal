<template>
  <div class="dashboard-container">
    <div class="dashboard-text">
      name:{{ name }}
    </div>
    <div class="dashboard-text">
      displayName:{{ displayName }}
    </div>
    <div class="dashboard-text">
      admin:{{ admin }}
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { UserModule } from '@/store/modules/user'

@Component({
  name: 'Dashboard'
})
export default class extends Vue {
  get name() {
    return UserModule.name
  }

  get displayName() {
    return UserModule.displayName
  }

  get admin() {
    return UserModule.admin
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }

  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
