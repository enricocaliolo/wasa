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

            conversationStore.init();
    
            switch(wsData.type) {
                case 'new_message':
                    const newMessage = new Message(wsData.payload.message);

                    const targetConversation = conversationStore.conversations.find(
                        conv => conv.conversationId === newMessage.conversationId
                    );
                    
                    if (targetConversation) {
                        targetConversation.lastMessage = newMessage;
                        
                        if (conversationStore.currentConversation?.conversationId === newMessage.conversationId 
                            || newMessage.isForwarded) {
                            targetConversation.messages.push(newMessage);
                        }
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
                    
                    if (!newConversation.messages) {
                        newConversation.messages = [];
                    }
                    
                    conversationStore.addConversation(newConversation);
                    break;
            }
        } catch (error) {
            console.error('Error processing WebSocket message:', error);
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
            reconnectInterval.value = Math.min(reconnectInterval.value * 2, maxReconnectInterval);
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

        ws.value.onerror = (error) => {
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