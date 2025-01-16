<script setup>
const props = defineProps({
	conversation: Object,
});
console.log(props.conversation);

const checkPhoto = () => {
	if (props.conversation.isGroup && props.conversation.photo) {
		return true;
	} else if (!props.conversation.isGroup) {
		const participant = props.conversation.participants.find(
			(participant) =>
				participant.username === props.conversation.name &&
				participant.icon !== undefined
		);
		if (participant) {
			return true;
		}
		return false;
	}
};

const getPhoto = () => {
	if (props.conversation.isGroup && props.conversation.photo) {
		return props.conversation.displayPhoto;
	} else if (!props.conversation.isGroup) {
		const participant = props.conversation.participants.find(
			(participant) => participant.username === props.conversation.name
		);
		const ret = participant ? participant.displayIcon : '';
		return ret;
	}
};
</script>

<template>
	<div class="avatar-container">
		<img v-if="checkPhoto()" :src="getPhoto()" class="avatar-image" />
		<span v-else class="avatar-placeholder">
			{{ props.conversation.name.charAt(0).toUpperCase() }}
		</span>
	</div>
</template>

<style>
.avatar-container {
	width: 50px;
	height: 50px;
	position: relative;
}

.avatar-image {
	width: 100%;
	height: 100%;
	border-radius: 50%;
	object-fit: cover;
}

.avatar-placeholder {
	width: 40px;
	height: 40px;
	background-color: #e5e7eb;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: 500;
	color: #374151;
}
</style>
