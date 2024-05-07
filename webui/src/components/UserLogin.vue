<script setup>
import Card from 'primevue/card';
import Password from 'primevue/password';
import InputGroup from 'primevue/inputgroup';
import InputText from "primevue/inputtext";
import InputGroupAddon from 'primevue/inputgroupaddon';
import Button from 'primevue/button';


import { ref } from 'vue';
const value = ref(null);
</script>

<template>
  <Card>
    <template #title>Log in</template>
    <template #content>
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-user"></i>
        </InputGroupAddon>
        <InputText placeholder="Username" />
      </InputGroup>
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-lock"></i>
        </InputGroupAddon>
        <Password v-model="value" placeholder="Password" :feedback="false" toggleMask />
      </InputGroup>
    </template>
    <template #footer>
      <div class="flex gap-3 mt-1">
        <Button label="Log in" class="w-full" />
      </div>
    </template>
  </Card>
</template>

