<script setup>
import Vertical from "@/components/legos/Vertical.vue";
import AppTopbar from "@/layout/AppTopbar.vue";
</script>

<template>
  <vertical :center-content="false">

    <template v-slot:header>
      <AppTopbar />
    </template>

    <template v-slot:main >
      <div>
        Main banana
      </div>
    </template>

    <template v-slot:footer>
      footer
    </template>

  </vertical>

</template>
