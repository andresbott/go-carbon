<script setup>
import Vertical from "@/components/legos/Vertical.vue";
import AppTopbar from "@/layout/AppTopbar.vue";
</script>

<template>
  <vertical>

    <template v-slot:header>
      <AppTopbar />
    </template>

    <template v-slot:default>
      Main
    </template>

    <template v-slot:footer>
      footer
    </template>

  </vertical>

</template>
