<script setup>
import Avatar from 'primevue/avatar'
import Drawer from 'primevue/drawer'
import Button from 'primevue/button'
const user = useUserStore()
const tasks = useTaskStore()
import { ref } from 'vue'
import { useUserStore } from '@/stores/user.js'
import { useTaskStore } from '@/stores/taskStore.js'
const visible = ref(false)

const logOut = () => {
    tasks.reset()
    user.logout()
}
</script>
<template>
    <Avatar
        icon="pi pi-user"
        class="mr-2"
        size="large"
        @click="visible = true"
        style="background-color: #ece9fc; color: #2a1261; cursor: pointer"
    />

    <Drawer v-model:visible="visible" :header="user.userName" style="width: 25rem" position="right">
        <Avatar
            icon="pi pi-user"
            class="mr-3"
            size="xlarge"
            @click="visible = true"
            style="background-color: #ece9fc; color: #2a1261"
        />

        <p>TODO make to look more like github</p>
        <p><Button label="Settings" icon="pi pi-cog" /></p>
        <p>
            <Button label="Logout" severity="danger" icon="pi pi-sign-out" @click="logOut()" />
        </p>
    </Drawer>
</template>
