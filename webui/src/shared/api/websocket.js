// useWebSocket.js
import { ref, onUnmounted, onMounted } from 'vue';
import { useConversationStore } from '../stores/conversation_store';
import { Message } from '../../modules/message/models/message';
import { useUserStore } from '../stores/user_store';
import { Conversation } from '../../modules/conversation/models/conversation';
import { User } from '../../modules/auth/models/user';
import { Reaction } from '../../modules/message/models/reaction';

export function useWebSocket() {
	const ws = ref(null);
	const isConnected = ref(false);
	const reconnectAttempts = ref(0);
	const maxReconnectAttempts = 5;
	const reconnectInterval = ref(1000);
	const maxReconnectInterval = 30000;
	const reconnectTimeout = ref(null);
	const userStore = useUserStore();

	const handleWebSocketMessage = (event) => {
		try {
			const wsData = JSON.parse(event.data);
			const conversationStore = useConversationStore();

			conversationStore.init();

			switch (wsData.type) {
				case 'new_message': {
					const newMessage = new Message(wsData.payload.message);

					const targetConversation = conversationStore.conversations.find(
						(conv) => conv.conversationId === newMessage.conversationId
					);

					if (targetConversation) {
						targetConversation.messages.push(newMessage);
					}
					break;
				}

				case 'new_conversation': {
					if (!wsData.payload || !wsData.payload.conversation) {
						console.error('Invalid conversation payload:', wsData);
						return;
					}

					const newConversation = new Conversation(wsData.payload.conversation);

					if (!newConversation.isGroup) {
						const otherParticipant = newConversation.participants.find(
							(participant) => participant.userId !== userStore.user.userId
						);
						if (otherParticipant) {
							newConversation.name = otherParticipant.username;
						}
					}

					if (!newConversation.messages) {
						newConversation.messages = [];
					}

					const existingConversation = conversationStore.conversations.find(
						(conv) => conv.conversationId === newConversation.conversationId
					);

					if (existingConversation) {
						existingConversation.participants =
							newConversation.participants.map(
								(participant) => new User(participant)
							);
						existingConversation.name = newConversation.name;
						existingConversation.photo = newConversation.photo;
						existingConversation.isGroup = newConversation.isGroup;
					} else {
						conversationStore.addConversation(newConversation);
					}
					break;
				}

				case 'messages_seen': {
					const { user_id, message_ids } = wsData.payload;

					if (!Array.isArray(message_ids)) {
						console.error('Invalid message_ids:', message_ids);
						return;
					}

					const targConv = conversationStore.conversations.find(
						(conv) => conv.conversationId === wsData.conversation_id
					);

					if (targConv) {
						message_ids.forEach((messageId) => {
							const message = targConv.messages.find(
								(m) => m.messageId === messageId
							);
							if (message && !message.isSeenBy(user_id)) {
								message.addSeenBy(user_id);
							}
						});
					}
					break;
				}

				case 'message_deletion': {
					const deletedMessage = new Message(wsData.payload.message);
					const targetConversation = conversationStore.conversations.find(
						(conv) => conv.conversationId === deletedMessage.conversationId
					);
					targetConversation.messages = targetConversation.messages.map(msg => {
						if (msg.messageId === deletedMessage.messageId) {
							return deletedMessage;
						}
						return msg;
					});
				}

				case 'reaction_update': {
					const targetConversation = conversationStore.conversations.find(
						(conv) => conv.conversationId === wsData.conversation_id
					);
				
					if (targetConversation) {
						const messageId = wsData.payload.reaction.message_id;
						const messageToUpdate = targetConversation.messages.find(
							(msg) => msg.messageId === messageId
						);
				
						if (messageToUpdate) {
							if (!messageToUpdate.reactions) {
								messageToUpdate.reactions = [];
							}
							messageToUpdate.reactions.push(new Reaction(wsData.payload.reaction));
						}
					}
					break;
				}

				case 'reaction_deletion': {
					const { reaction } = wsData.payload;
					
					const targetConversation = conversationStore.conversations.find(
						(conv) => conv.conversationId === wsData.conversation_id
					);
				
					if (targetConversation) {
						const targetMessage = targetConversation.messages.find(
							(msg) => msg.messageId === reaction.message_id
						);
						if (targetMessage) {
							targetMessage.reactions = targetMessage.reactions.filter(
								(r) => r.reactionId !== reaction.reaction_id
							);
						}
					}
					break;
				}
			}
		} catch (error) {
			console.error('Error processing WebSocket message:', error);
		}
	};

	const sendMessagesSeen = (conversationId, messageIds) => {
		if (!Array.isArray(messageIds) || messageIds.length === 0) {
			console.warn('Invalid messageIds:', messageIds);
			return;
		}

		if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
			console.warn('WebSocket not connected');
			return;
		}

		try {
			const message = {
				type: 'messages_seen',
				conversation_id: conversationId,
				payload: {
					message_ids: messageIds,
					user_id: userStore.user.userId,
				},
				timestamp: new Date(),
			};
			ws.value.send(JSON.stringify(message));
		} catch (e) {
			console.log(e);
		}
	};

	const attemptReconnect = () => {
		if (reconnectAttempts.value >= maxReconnectAttempts) {
			return;
		}

		if (reconnectTimeout.value) {
			clearTimeout(reconnectTimeout.value);
		}

		reconnectTimeout.value = setTimeout(() => {
			reconnectAttempts.value++;
			connectWebSocket();
			reconnectInterval.value = Math.min(
				reconnectInterval.value * 2,
				maxReconnectInterval
			);
		}, reconnectInterval.value);
	};

	const connectWebSocket = () => {
		if (!userStore.user?.userId) {
			return;
		}

		const wsUrl = `ws://localhost:3000/ws?token=${userStore.user.userId}`;

		if (ws.value && ws.value.readyState === WebSocket.OPEN) {
			ws.value.close();
		}

		ws.value = new WebSocket(wsUrl);

		const pingInterval = setInterval(() => {
			if (ws.value && ws.value.readyState === WebSocket.OPEN) {
				ws.value.send(JSON.stringify({ type: 'ping' }));
			}
		}, 30000);

		ws.value.onopen = () => {
			isConnected.value = true;
			reconnectAttempts.value = 0;
			reconnectInterval.value = 1000;
		};

		ws.value.onmessage = handleWebSocketMessage;

		ws.value.onclose = (event) => {
			isConnected.value = false;
			clearInterval(pingInterval);
			if (event.code !== 1000) {
				attemptReconnect();
			}
		};

		ws.value.onerror = () => {
			isConnected.value = false;
		};
	};

	const disconnect = () => {
		if (reconnectTimeout.value) {
			clearTimeout(reconnectTimeout.value);
			reconnectTimeout.value = null;
		}

		if (ws.value) {
			ws.value.close(1000, 'Intentional disconnection');
			ws.value = null;
		}
		reconnectAttempts.value = 0;
		reconnectInterval.value = 1000;
	};

	onMounted(() => {
		connectWebSocket();
	});

	onUnmounted(() => {
		disconnect();
	});

	return {
		isConnected,
		connect: connectWebSocket,
		disconnect,
		sendMessagesSeen,
	};
}
