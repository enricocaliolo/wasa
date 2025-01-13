// useWebSocket.js
import { ref, onUnmounted } from 'vue';
import { useConversationStore } from '../stores/conversation_store';
import { Message } from '../../modules/message/models/message';
import { useUserStore } from '../stores/user_store';
import { Conversation } from '../../modules/conversation/models/conversation';

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

            switch(wsData.type) {
                case 'new_message':
                    const newMessage = new Message(wsData.payload.message);
        
                    if (conversationStore.currentConversation?.conversationId === newMessage.conversationId) {
                        conversationStore.currentConversation.messages.push(newMessage);
                    }
                    
                    const affectedConversation = conversationStore.conversations.find(
                        conv => conv.conversationId === newMessage.conversationId
                    );
                    
                    if (affectedConversation) {
                        affectedConversation.lastMessage = newMessage;
                    }

                    if (!(conversationStore.currentConversation?.conversationId === newMessage.conversationId) && newMessage.isForwarded) {
                        const conversation = conversationStore.conversations.find((c) => c.conversationId === newMessage.conversationId);
                        conversation?.messages.push(newMessage);
                    }
                    break;

                case 'new_conversation':
                    const newConversation = new Conversation(wsData.payload.conversation);

                    if (!newConversation.isGroup) {
                        const otherParticipant = newConversation.participants.find(
                            (participant) => participant.userId !== userStore.user.userId
                        );
                        if (otherParticipant) {
                            newConversation.name = otherParticipant.username;
                        }
                    }
                    conversationStore.init();
                    conversationStore.addConversation(newConversation);
                    break;
            }
        } catch (error) {
            console.error('Error processing message:', error);
        }
    };

    const attemptReconnect = () => {
        if (reconnectAttempts.value >= maxReconnectAttempts) {
            console.log('Max reconnection attempts reached');
            return;
        }

        if (reconnectTimeout.value) {
            clearTimeout(reconnectTimeout.value);
        }

        reconnectTimeout.value = setTimeout(() => {
            console.log(`Attempting to reconnect (${reconnectAttempts.value + 1}/${maxReconnectAttempts})`);
            reconnectAttempts.value++;
            connectWebSocket();
            reconnectInterval.value = Math.min(reconnectInterval.value * 2, maxReconnectInterval);
        }, reconnectInterval.value);
    };

    const connectWebSocket = () => {
        if (!userStore.user?.userId) {
            console.log('No user logged in, skipping WebSocket connection');
            return;
        }

        console.log('Connecting to WebSocket...');
        const wsUrl = `ws://localhost:3000/ws?token=${userStore.user.userId}`;

        if (ws.value && ws.value.readyState === WebSocket.OPEN) {
            ws.value.close();
        }
        
        ws.value = new WebSocket(wsUrl);
        
        const pingInterval = setInterval(() => {
            if (ws.value && ws.value.readyState === WebSocket.OPEN) {
                ws.value.send(JSON.stringify({ type: 'ping' }));
            }
        }, 30000); // Send ping every 30 seconds

        ws.value.onopen = () => {
            console.log('WebSocket connected');
            isConnected.value = true;
            reconnectAttempts.value = 0;
            reconnectInterval.value = 1000;
        };

        ws.value.onmessage = handleWebSocketMessage;

        ws.value.onclose = (event) => {
            console.log('WebSocket disconnected', event.code, event.reason);
            isConnected.value = false;
            clearInterval(pingInterval);
            if (event.code !== 1000) {
                attemptReconnect();
            }
        };

        ws.value.onerror = (error) => {
            console.error('WebSocket error:', error);
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

    onUnmounted(() => {
        disconnect();
    });

    return {
        isConnected,
        connect: connectWebSocket,
        disconnect
    };
}