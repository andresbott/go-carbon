<script setup>
import Card from 'primevue/card'
import Password from 'primevue/password'
import InputGroup from 'primevue/inputgroup'
import InputText from 'primevue/inputtext'
import InputGroupAddon from 'primevue/inputgroupaddon'
import Button from 'primevue/button'
import Message from 'primevue/message'

import { computed, onBeforeMount, onMounted } from 'vue'

import { ref } from 'vue'
import { useUserStore } from '@/stores/user.js'
import router from '@/router/index.js'
import LoadingScreen from '@/components/loadingScreen.vue'

const user = useUserStore()

const userRef = ref(null)
const passRef = ref(null)

const load = () => {
    user.login(userRef.value, passRef.value)
}
const status = ref('out')
const visible = ref(true)
</script>

<template>
    <Card>
        <template #title>Log in</template>
        <template #content>
            <div v-focustrap class="flex flex-column items-center gap-4">
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
                    <Password
                        v-model="passRef"
                        v-on:keyup.enter="load"
                        placeholder="Password"
                        :feedback="false"
                        toggleMask
                    />
                </InputGroup>

                <Message v-if="user.wrongPw" severity="error" closable
                    >Wrong username or password</Message
                >
                <Button label="Log in" class="w-full" @click="load" />
            </div>
        </template>
    </Card>
    <loadingScreen v-if="user.loading" />
</template>
