<script setup>
import { useConversationStore } from '@/shared/stores/conversation_store';
import SidebarComponent from '../components/SidebarComponent.vue';
import CurrentConversationView from './CurrentConversationView.vue';

const currentConversationStore = useConversationStore();
</script>

<template>
	<main>
		<SidebarComponent />
		<CurrentConversationView
			v-if="currentConversationStore.currentConversation"
			:conversation="currentConversationStore.currentConversation"
		/>
		<div v-else class="current-conversation"></div>
	</main>
</template>

<style scoped>

main {
	height: 100vh;
	width: 90vw;
	margin: 0 auto;
	max-width: 1500px;
	display: grid;
	grid-template-columns: 1fr 2fr;
	padding: 5vh 0;
}

.current-conversation {
	background-color: grey;
}
</style>
