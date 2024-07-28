<script setup>
import Vertical from '@/components/legos/Vertical.vue'
import Horizontal from '@/components/legos/horizontal.vue'
import TopBar from '@/views/parts/topbar.vue'
import ProjectIcon from '@/components/icons/ProjectIcon.vue'
</script>

<template>
    <vertical :center-content="true">
        <template v-slot:header>
            <TopBar />
        </template>

        <template v-slot:main>
            <horizontal :center-content="true">
                <template v-slot:default>
                    <div>Landing page content</div>
                    <ProjectIcon />
                    <router-view />
                </template>
            </horizontal>
        </template>
    </vertical>
</template>
