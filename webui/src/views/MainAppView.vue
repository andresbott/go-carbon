<script setup>
import Vertical from '@/components/legos/Vertical.vue'
import Horizontal from '@/components/legos/horizontal.vue'
import TopBar from '@/views/parts/topbar.vue'
import {useUserStore} from "@/stores/user.js";

const user = useUserStore()
console.log(user.banana)
</script>

<template>
  <vertical :center-content="true">
    <template v-slot:header>
      <TopBar />
    </template>

    <template v-slot:main>
      <horizontal :center-content="true">
        <template v-slot:default>
          <div>Main banana</div>
          <router-view />
        </template>
      </horizontal>
    </template>

    <template v-slot:footer>my footer </template>
  </vertical>
</template>
