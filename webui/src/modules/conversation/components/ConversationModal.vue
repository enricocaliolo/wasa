<script setup>
import { ref, onMounted, computed } from "vue";
import { userAPI } from "@/modules/auth/api/user-api";
import { useConversationStore } from "@/shared/stores/conversation_store";
import { useUserStore } from "@/shared/stores/user_store";

const props = defineProps({
    show: Boolean,
});

const conversationStore = useConversationStore();
const userStore = useUserStore();

const emit = defineEmits(["close"]);

const groupName = ref("");
const currentUsersForConversation = ref([userStore.user]);
const allUsers = ref([]);

const availableUsers = computed(() => {
    return allUsers.value
        .filter(user => user.userId !== userStore.user.userId)
        .map(user => ({
            ...user,
            isSelected: currentUsersForConversation.value.some(
                selectedUser => selectedUser.userId === user.userId
            )
        }));
});

const isGroupChat = computed(() => currentUsersForConversation.value.length > 2);
const isCreateButtonDisabled = computed(() => 
    currentUsersForConversation.value.length < 2 || 
    (isGroupChat.value && !groupName.value)
);

onMounted(async () => {
    try {
        const users = await userAPI.getAllUsers();
        allUsers.value = users;
    } catch (error) {
        console.error('Error loading users:', error);
    }
});

function toggleUser(user) {
    const index = currentUsersForConversation.value.findIndex(
        u => u.userId === user.userId
    );
    
    if (index === -1) {
        currentUsersForConversation.value.push(user);
    } else {
        currentUsersForConversation.value = currentUsersForConversation.value
            .filter(u => u.userId !== user.userId);
    }
}

function closeModal() {
    groupName.value = "";
    currentUsersForConversation.value = [userStore.user];
    emit("close");
}

async function createConversation() {
    await conversationStore.createConversation({
        currentUsers: currentUsersForConversation.value,
        groupName: groupName.value,
    });
    closeModal();
}
</script>

<template>
    <BaseModal :show="show" title="Create Conversation" @close="closeModal">
        <div class="modal-content">
            <div class="selected-users-section">
                <h3 class="section-title">Selected Users ({{ currentUsersForConversation.length - 1 }})</h3>
                <div class="selected-users">
                    <div
                        v-for="user in currentUsersForConversation.slice(1)"
                        :key="user.userId"
                        class="selected-user-item"
                    >
                        {{ user.username }}
                    </div>
                </div>
            </div>

            <div v-if="isGroupChat" class="group-name-section">
                <h3 class="section-title">Group Name</h3>
                <input
                    type="text"
                    placeholder="Enter group name"
                    v-model="groupName"
                    class="group-name-input"
                />
            </div>

            <div class="users-list">
                <h3 class="section-title">Select Users</h3>
                <div class="users-grid">
                    <div
                        v-for="user in availableUsers"
                        :key="user.userId"
                        class="user-item"
                        :class="{ 'user-selected': user.isSelected }"
                        @click="toggleUser(user)"
                    >
						<span v-if="user.icon">
							<img :src="`data:image/jpeg;base64,${user.icon}`" alt="user icon" class="user-icon">
						</span>
						<div v-else class="default-icon">
                            {{ user.username.charAt(0).toUpperCase() }}
                        </div>

                        <span class="username">{{ user.username }}</span>
                        <span v-if="user.isSelected" class="selected-indicator">✓</span>
                    </div>
                </div>
            </div>
        </div>

        <template #footer>
            <div class="modal-footer">
                <button class="cancel-button" @click="closeModal">Cancel</button>
                <button
                    class="create-button"
                    @click="createConversation"
                    :disabled="isCreateButtonDisabled"
                >
                    Create
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
    min-width: 400px
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
    justify-content: space-between;
    align-items: center;
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
    color: #2196f3;
    font-weight: bold;
}

.group-name-section {
    margin: 1rem 0;
}

.group-name-input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 1rem;
}

.group-name-input:focus {
    border-color: #2196f3;
    outline: none;
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

.create-button {
    background-color: #2196f3;
    color: white;
    border: none;
}

.create-button:hover:not(:disabled) {
    background-color: #1976d2;
}

.create-button:disabled {
    background-color: #bbdefb;
    cursor: not-allowed;
}
</style>