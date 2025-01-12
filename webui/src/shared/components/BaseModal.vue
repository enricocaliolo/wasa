<script setup>
defineProps({
	show: Boolean,
	title: {
		type: String,
		default: "",
	},
});

const emit = defineEmits(["close"]);

function closeModal() {
	emit("close");
}
</script>

<template>
	<div v-if="show" class="modal-overlay">
		<div class="base-modal">
			<header class="modal-header">
				<h3 v-if="title">{{ title }}</h3>
				<button class="close-button" @click="closeModal">×</button>
			</header>
			<div class="modal-content">
				<slot></slot>
			</div>
			<footer class="modal-footer">
				<slot name="footer"></slot>
			</footer>
		</div>
	</div>
</template>

<style scoped>
.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background: rgba(0, 0, 0, 0.5);
	display: grid;
	place-items: center;
	z-index: 1000;
}

.base-modal {
	background-color: white;
	border-radius: 8px;
	padding: 1rem;
	min-width: 300px;
	max-width: 90vw;
	max-height: 90vh;
	overflow-y: auto;
}

.modal-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 1rem;
}

.close-button {
	background: none;
	border: none;
	font-size: 1.5rem;
	cursor: pointer;
	padding: 0.25rem 0.5rem;
}

.modal-content {
	margin-bottom: 1rem;
}

.modal-footer {
	display: flex;
	justify-content: flex-end;
	gap: 0.5rem;
}
</style>
