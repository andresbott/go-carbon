<script setup>
import Vertical from "@/components/legos/Vertical.vue";
import AppTopbar from "@/layout/AppTopbar.vue";
import Horizontal from "@/components/legos/horizontal.vue";
import UserLogin from "@/components/UserLogin.vue";
</script>

<template>
  <vertical :center-content="true">

    <template v-slot:header>
      <horizontal :center-content="true">
        <template v-slot:left>
          <router-link to="/" class="layout-topbar-logo">
            <span>SAKAI</span>
          </router-link>
        </template>

        <template v-slot:right>
          <button  class="">
            <i class="pi pi-user"></i>
            <span>Profile</span>
          </button>
        </template>
      </horizontal>
    </template>

    <template v-slot:main >
      <horizontal :center-content="true">
        <template v-slot:default>
          <UserLogin />
        </template>
      </horizontal>
    </template>

    <template v-slot:footer>
    </template>

  </vertical>

</template>
