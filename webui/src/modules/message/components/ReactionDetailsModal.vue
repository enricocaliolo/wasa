<script setup>
import { computed } from 'vue';
import BaseModal from '../../../shared/components/BaseModal.vue';

const props = defineProps({
  show: Boolean,
  reactions: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['close']);

const groupedReactions = computed(() => {
  return props.reactions.reduce((acc, reaction) => {
    if (!acc[reaction.reaction]) {
      acc[reaction.reaction] = [];
    }
    acc[reaction.reaction].push(reaction.user);
    return acc;
  }, {});
});

function closeModal() {
  emit('close');
}
</script>

<template>
  <BaseModal :show="show" title="Reactions" @close="closeModal">
    <div class="reactions-list">
      <div v-for="(users, emoji) in groupedReactions" :key="emoji" class="reaction-group">
        <div class="emoji-header">
          <span class="emoji">{{ emoji }}</span>
          <span class="count">{{ users.length }}</span>
        </div>
        <div class="users-list">
          <div v-for="user in users" :key="user.userId" class="user-item">
            <img 
              v-if="user.displayIcon" 
              :src="user.displayIcon" 
              :alt="user.username" 
              class="user-icon" 
            />
            <span class="username">{{ user.username }}</span>
          </div>
        </div>
      </div>
    </div>
  </BaseModal>
</template>

<style scoped>
.reactions-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  min-width: 300px;
}

.reaction-group {
  border-bottom: 1px solid #eee;
  padding-bottom: 1rem;
}

.reaction-group:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.emoji-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.emoji {
  font-size: 1.5rem;
}

.count {
  color: #666;
  font-size: 0.9rem;
}

.users-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem;
  border-radius: 4px;
  transition: background-color 0.2s ease;
}

.user-item:hover {
  background-color: #f5f5f5;
}

.user-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.username {
  font-size: 0.9rem;
}
</style>