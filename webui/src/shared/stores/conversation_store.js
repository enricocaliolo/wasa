import { ref } from "vue";
import { defineStore } from "pinia";
import { messagesAPI } from "@/modules/message/api/message_api";
import { useUserStore } from "./user_store";
import { Message } from "../../modules/message/models/message";
import { Reaction } from "../../modules/message/models/reaction";
import { conversationAPI } from "../../modules/conversation/api/conversation-api";
import { imageConverter } from "../../modules/message/helper/image_converter";

export const useConversationStore = defineStore("conversationStore", () => {
	const userStore = useUserStore();

	const conversations = ref([]);
	const currentConversation = ref();
	const replyMessage = ref(null);
	const showGroupDetails = ref(false)

	function init() {
        if (!conversations.value) {
            conversations.value = [];
        }
    }

	function setCurrentConversation(conversation) {
		currentConversation.value = conversation;
	}

	async function getUserConversations() {
		conversations.value = await conversationAPI.getUserConversations();
	}
	
	async function createConversation({currentUsers, groupName}) {
		if (currentUsers.length > 2 && !groupName) {
			alert("Please, insert a group name");
			return;
		} else if (currentUsers.length === 2) {
			groupName = currentUsers[1].username;
		}
	
		const conversation = await conversationAPI.createConversation(
			currentUsers.map((user) => user.userId),
			groupName,
		);
	
		if (currentUsers.length === 2) {
			conversation.name = currentUsers[1].username;
		}
	}

	function addConversation(conv) {
		const existingConv = conversations.value.find(
			c => c.conversationId === conv.conversationId
		);
		
		if (!existingConv) {
			conversations.value.push(conv);
			
			if (conversations.value.length === 1) {
				setCurrentConversation(conv);
			}
		}
	}

	function setReplyMessage(message) {
		replyMessage.value = message;
	}

	async function toggleGroupDetails(value) {
		showGroupDetails.value = value
	}

	async function addReaction(message_id, _reaction) {
		const data = await messagesAPI.commentMessage(
			currentConversation.value.conversationId,
			message_id,
			_reaction,
		);
		const reaction = Reaction.fromJSON(data);

		const messageToUpdate = currentConversation.value.messages.find(
			(message) => message.messageId === message_id,
		);

		if (messageToUpdate) {
			if (!messageToUpdate.reactions) {
				messageToUpdate.reactions = [];
			}
			messageToUpdate.reactions.push(reaction);
		}

		return reaction;
	}

	async function deleteReaction(message_id, reaction_id) {
		await messagesAPI.uncommentMessage(
			currentConversation.value.conversationId,
			message_id,
			reaction_id,
		);

		for (const message of currentConversation.value.messages) {
			if (message.messageId === message_id) {
				message.reactions = message.reactions.filter(
					(reaction) => reaction.reactionId !== reaction_id,
				);
			}
		}

		return;
	}

	async function updateGroupName(name) {
		await conversationAPI.updateGroupName(
			currentConversation.value.conversationId,
			name,
		);
		currentConversation.value.name = name
		return
	}

	async function updateGroupPhoto(photo) {
		await conversationAPI.updateGroupPhoto(
			currentConversation.value.conversationId,
			photo,
		)
		currentConversation.value.photo = await imageConverter.fileToBase64(photo)
		return
	}

	async function leaveGroup() {
		try{
			await conversationAPI.leaveGroup(currentConversation.value.conversationId)
			conversations.value = conversations.value.filter(
				(conv) => conv.conversationId !== currentConversation.value.conversationId
			)
			currentConversation.value = null
		} catch(e) {
			throw e
		}
	}


	async function sendMessage({
		content,
		content_type = "text",
		replied_to_message = null,
		source_conversation_id = null,
		destination_conversation_id = null,
	}) {
		const conversationId =
			source_conversation_id || currentConversation.value.conversationId;
		replied_to_message =
			replied_to_message !== null ? replyMessage.value.messageId : null;

		const data = await messagesAPI.sendMessage({
			conversation_id: conversationId,
			content,
			content_type,
			replied_to_message,
			destination_conversation_id,
		});

		const message = Message.fromJSON(data);
		message.sender = userStore.getUser();

		return message;
	}

	return {
		conversations,
		currentConversation,
		replyMessage,
		setCurrentConversation,
		sendMessage,
		addConversation,
		setReplyMessage,
		addReaction,
		deleteReaction,
		updateGroupName,
		updateGroupPhoto,
		toggleGroupDetails,
		showGroupDetails,
		leaveGroup,
		createConversation,
		getUserConversations,
		init
	};
});
