<script setup>
import Horizontal from "@/components/legos/horizontal.vue";
</script>

<template>
  <horizontal :center-content="true">
    <template v-slot:left>
      <router-link to="/" class="layout-topbar-logo">
        <span>SAKAI</span>
      </router-link>
    </template>

    <template v-slot:default>
      Search
    </template>

    <template v-slot:right>
      <button  class="">
        <i class="pi pi-user"></i>
        <span>Profile</span>
      </button>
    </template>
  </horizontal>

</template>

<style lang="scss" scoped></style>
