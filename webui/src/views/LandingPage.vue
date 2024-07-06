<script setup>
import Vertical from '@/components/legos/Vertical.vue'
import Horizontal from '@/components/legos/horizontal.vue'

import {useUserStore} from "@/stores/user.js";
import Avatar from "primevue/avatar";
import InputText from "primevue/inputtext";
import UserProfile from "@/components/UserProfile.vue";
import Button from "primevue/button";
import TopBar from "@/views/parts/topbar.vue";


</script>

<template>
  <vertical :center-content="true">
    <template v-slot:header>


          <TopBar />


    </template>

    <template v-slot:main>
      <horizontal :center-content="true">
        <template v-slot:default>
          <div>Landing page content</div>
          <router-view />
        </template>
      </horizontal>
    </template>


  </vertical>
</template>
