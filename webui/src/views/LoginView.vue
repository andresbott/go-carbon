<script setup>
import Vertical from '@/components/legos/Vertical.vue'
import Horizontal from '@/components/legos/horizontal.vue'
import UserLogin from '@/components/UserLogin.vue'
</script>

<template>
    <vertical :center-content="true">
        <template v-slot:main>
            <horizontal :center-content="true">
                <template v-slot:default>
                    <UserLogin />
                </template>
            </horizontal>
        </template>

        <template v-slot:footer> </template>
    </vertical>
</template>
