import { defineStore } from 'pinia';
import { User } from '../../modules/auth/models/user';
import { ref } from 'vue';
import { userAPI } from '../../modules/auth/api/user-api';
import { imageConverter } from '../../modules/message/helper/image_converter';
import router from '../router/router';

export const useUserStore = defineStore('userStore', () => {
	const user = ref(User.fromJSON({}));
	function setUser(_user) {
		user.value = _user;
	}

	function getUser() {
		return user.value;
	}

	async function updateUsername(_username) {
		const response = await userAPI.changeUsername(_username);
		if (response) {
			user.value.username = _username;
		}
	}

	async function updateIcon(_icon) {
		await userAPI.changeIcon(_icon);
		user.value.icon = await imageConverter.fileToBase64(_icon);
		return;
	}

	async function logout() {
		setUser(User.fromJSON({}));
		localStorage.removeItem('username');
		router.push({path: '/login', replace: true,});
	}

	return { user, setUser, getUser, updateUsername, updateIcon, logout, };
});
