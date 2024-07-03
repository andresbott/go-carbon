<script setup>
import Card from 'primevue/card';
import Password from 'primevue/password';
import InputGroup from 'primevue/inputgroup';
import InputText from "primevue/inputtext";
import InputGroupAddon from 'primevue/inputgroupaddon';
import Button from 'primevue/button';


import { ref } from 'vue';
import {useUserStore} from "@/stores/user.js";
const value = ref(null);

const user = useUserStore()


const userRef = ref(null)
const passRef = ref(null)


const load = () => {
  user.login(userRef.value,passRef.value)
};



</script>

<template>
  <Card>
    <template #title>Log in</template>
    <template #content>
      <div class="flex flex-column items-center gap-4">

      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-user"></i>
        </InputGroupAddon>
        <InputText placeholder="Username" v-on:keyup.enter="load" v-model="userRef" />
      </InputGroup>

      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-lock"></i>
        </InputGroupAddon>
        <Password v-model="passRef" v-on:keyup.enter="load" placeholder="Password" :feedback="false" toggleMask />
      </InputGroup>
        <Button label="Log in" class="w-full" @click="load" />
      </div>
    </template>
  </Card>
</template>

