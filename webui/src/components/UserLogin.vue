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
// user.bla()
// user.login()

const userRef = ref(null)
const passRef = ref(null)


const load = () => {

  user.bla(userRef.value,passRef.value)
};



</script>

<template>
  <Card>
    <template #title>Log in</template>
    <template #content>
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-user"></i>
        </InputGroupAddon>
        <InputText placeholder="Username" v-model="userRef" />
      </InputGroup>
      <InputGroup>
        <InputGroupAddon>
          <i class="pi pi-lock"></i>
        </InputGroupAddon>
        <Password v-model="passRef" placeholder="Password" :feedback="false" toggleMask />
      </InputGroup>
    </template>
    <template #footer>
      <div class="flex gap-3 mt-1">
        <Button label="Log in" class="w-full" @click="load" />
      </div>
    </template>
  </Card>
</template>

