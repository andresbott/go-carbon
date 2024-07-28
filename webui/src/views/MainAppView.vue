<script setup>
import Vertical from '@/components/legos/Vertical.vue'
import Horizontal from '@/components/legos/horizontal.vue'
import TopBar from '@/views/parts/topbar.vue'
import { useUserStore } from '@/stores/user.js'
import TaskList from '@/views/parts/TaskList.vue'

const user = useUserStore()
</script>

<template>
    <vertical :center-content="false">
        <template v-slot:header>
            <TopBar />
        </template>

        <template v-slot:main>
            <horizontal :center-content="true">
                <template v-slot:default>
                    <TaskList />
                    <router-view />
                </template>
            </horizontal>
        </template>

        <template v-slot:footer>my footer </template>
    </vertical>
</template>
