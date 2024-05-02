<template>
  <div class="about">
    <h1>This is an about page</h1>
  </div>

  <vertical>
    <template v-slot:header>
     head
    </template>

Main content


    <template v-slot:footer>
      footer
    </template>

  </vertical>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
<script setup>
import Vertical from "@/components/legos/Vertical.vue";
</script>