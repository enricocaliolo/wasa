<script setup>
import InputComponent from "@/modules/conversation/components/InputComponent.vue";
import MessageComponent from "../../message/components/MessageComponent.vue";
import { useConversationStore } from "../../../shared/stores/conversation_store";
import {computed} from 'vue'
import GroupDetails from "../components/GroupDetails.vue";
import IconConversation from "../../../shared/components/IconConversation.vue";

const props = defineProps({
	conversation: Object,
});

const conversationStore = useConversationStore();

function changeGroupName() {
	conversationStore.toggleGroupDetails(!conversationStore.showGroupDetails);
}

</script>

<template>
	<div class="conversation-wrapper">
		<header @click="changeGroupName">
			<IconConversation :conversation="conversationStore.currentConversation" />
			<h1>{{ conversation.name }}</h1>
		</header>
		<GroupDetails v-if="conversationStore.showGroupDetails" :conversation="conversation" />
		<div v-else class="messages-box">	
			<MessageComponent
				v-for="message in conversation.messages"
				:key="message.messageId"
				:message="message"
			/>
		</div>
		<InputComponent v-if="!conversationStore.showGroupDetails" />
	</div>
</template>

<style scoped>
header {
    padding: 1em;
    background-color: grey;
    display: flex;
    align-items: center;
    gap: 1em;
    cursor: pointer;
}

.conversation-wrapper {
    display: grid;
    grid-template-rows: auto 1fr auto;
    min-height: 0;
}

.messages-box {
    background-color: lightblue;
    overflow-y: auto;
    padding: 1em;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    min-height: 0;
    gap: 10px;
}
</style>
