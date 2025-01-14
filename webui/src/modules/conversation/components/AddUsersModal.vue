<script setup>
import { ref, onMounted, computed } from "vue";
import { userAPI } from "@/modules/auth/api/user-api";
import { useConversationStore } from "@/shared/stores/conversation_store";
import { useUserStore } from "@/shared/stores/user_store";
import { User } from "../../auth/models/user";

const props = defineProps({
    show: Boolean,
    conversation: Object
});

const conversationStore = useConversationStore();

const emit = defineEmits(["close"]);
const selectedUsers = ref([]);
const allUsers = ref([]);

const availableUsers = computed(() => {
    const existingUserIds = props.conversation.participants.map(p => p.userId);
    
    return allUsers.value
        .filter(user => !existingUserIds.includes(user.userId))
        .map(user => {
            return {
                user: user,
                isSelected: selectedUsers.value.some(
                    selectedUser => selectedUser.userId === user.userId
                )
            };
        });
});

onMounted(async () => {
    try {
        const users = await userAPI.getAllUsers();
        allUsers.value = users;
        console.log('Loaded users:', users);
    } catch (error) {
        console.error('Error loading users:', error);
    }
});

function toggleUser(user) {
    const index = selectedUsers.value.findIndex(
        u => u.userId === user.userId
    );
    
    if (index === -1) {
        selectedUsers.value.push(user);
    } else {
        selectedUsers.value = selectedUsers.value
            .filter(u => u.userId !== user.userId);
    }
}

function closeModal() {
    selectedUsers.value = [];
    emit("close");
}

async function addUsers() {
    if (selectedUsers.value.length > 0) {
        await conversationStore.addGroupMembers({
            users: selectedUsers.value,
            conversationId: props.conversation.conversationId
        });
        closeModal();
    }
}
</script>

<template>
    <BaseModal :show="show" title="Add Users to Group" @close="closeModal">
        <div class="modal-content">
            <div class="users-list">
                <h3 class="section-title">Select Users to Add</h3>
                <div v-if="availableUsers.length === 0" class="no-users">
                    No more users available to add
                </div>
                <div v-else class="users-grid">
                    <div
                        v-for="userItem in availableUsers"
                        :key="userItem.user.userId"
                        class="user-item"
                        :class="{ 'user-selected': userItem.isSelected }"
                        @click="toggleUser(userItem.user)"
                    >
                        <span v-if="userItem.user.icon">
                            <img :src="`data:image/jpeg;base64,${userItem.user.icon}`" alt="user icon" class="user-icon">
                        </span>
                        <div v-else class="default-icon">
                            {{ userItem.user.username.charAt(0).toUpperCase() }}
                        </div>

                        <span class="username">{{ userItem.user.username }}</span>
                        <span v-if="userItem.isSelected" class="selected-indicator">✓</span>
                    </div>
                </div>
            </div>

            <div v-if="selectedUsers.length > 0" class="selected-users-section">
                <h3 class="section-title">Selected Users ({{ selectedUsers.length }})</h3>
                <div class="selected-users">
                    <div
                        v-for="user in selectedUsers"
                        :key="user.userId"
                        class="selected-user-item"
                    >
                        {{ user.username }}
                    </div>
                </div>
            </div>
        </div>

        <template #footer>
            <div class="modal-footer">
                <button class="cancel-button" @click="closeModal">Cancel</button>
                <button
                    class="add-button"
                    @click="addUsers"
                    :disabled="selectedUsers.length === 0"
                >
                    Add to Group
                </button>
            </div>
        </template>
    </BaseModal>
</template>

<style scoped>
.modal-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    max-height: 70vh;
    overflow-y: auto;
    padding: 1rem;
}

.no-users {
    text-align: center;
    padding: 2rem;
    color: #666;
    font-style: italic;
}

.section-title {
    font-size: 1rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: #333;
}

.users-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 0.5rem;
}

.user-icon {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
}

.default-icon {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background-color: #2196f3;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
}

.user-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem;
    background-color: #f5f5f5;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;
}

.user-item:hover {
    background-color: #e0e0e0;
}

.user-selected {
    background-color: #e3f2fd;
    border: 1px solid #2196f3;
}

.selected-indicator {
    margin-left: auto;
    color: #2196f3;
    font-weight: bold;
}

.selected-users {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
}

.selected-user-item {
    background-color: #e3f2fd;
    color: #1976d2;
    padding: 0.5rem 1rem;
    border-radius: 20px;
    font-size: 0.875rem;
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    padding-top: 1rem;
}

button {
    padding: 0.5rem 1rem;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;
}

.cancel-button {
    background-color: #f5f5f5;
    border: 1px solid #ddd;
}

.cancel-button:hover {
    background-color: #e0e0e0;
}

.add-button {
    background-color: #2196f3;
    color: white;
    border: none;
}

.add-button:hover:not(:disabled) {
    background-color: #1976d2;
}

.add-button:disabled {
    background-color: #bbdefb;
    cursor: not-allowed;
}
</style>