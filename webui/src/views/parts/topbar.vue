<script setup>
import Horizontal from '@/components/legos/horizontal.vue'
import InputText from 'primevue/inputtext'
import UserProfile from '@/components/UserProfile.vue'
import { useUserStore } from '@/stores/user.js'
import Button from 'primevue/button'
import ProjectIcon from '@/components/icons/ProjectIcon.vue'
const user = useUserStore()
</script>

<template>
    <horizontal :center-content="true">
        <template v-slot:left>
            <router-link to="/app" class="layout-topbar-logo">
                <ProjectIcon />
            </router-link>
        </template>

<!--        <InputText placeholder="Search" type="text" class="w-32 sm:w-auto" />-->
        <template v-slot:right>
            <UserProfile v-if="user.isLoggedIn" />
            <router-link v-if="!user.isLoggedIn" to="/login" class="layout-topbar-logo">
                <Button label="Login" icon="pi pi-sign-in" />
            </router-link>
        </template>
    </horizontal>
</template>
