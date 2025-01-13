import { useUserStore } from "@/shared/stores/user_store";
import api from "../../../shared/api/api";
import { User } from "../models/user";
import { useWebSocket } from "../../../shared/api/websocket";


export const userAPI = {
	login: async (_username) => {
		try {
			const userStore = useUserStore();
			const response = await api.put("/session", {
				username: _username,
			});

			if (response.data) {
				userStore.setUser(User.fromJSON(response.data));
				localStorage.setItem("username", _username)
			}

			// const {connect, disconnect } = useWebSocket();

			// connect();

			// const unwatch = import.meta.hot?.on('vite:beforeUpdate', () => {
			// 	disconnect();
			// });


			return true;
		} catch (error) {
			throw error;
		}
	},
	findUser: async (_username) => {
		try {
			const response = await api.get("/users/search", {
				params: {
					username: _username,
				},
			});

			const user = User.fromJSON(response.data);
			user.username = _username;

			return user;
		} catch (error) {
			throw error;
		}
	},
	changeUsername: async (_username) => {
		try {
			const response = await api.put("/settings/profile/username", {
				username: _username,
			});
			const user = User.fromJSON(response.data);
			return user;
		} catch (e) {
			if (e.response.status === 409) {
				return new Error("Username already exists! Choose another.")
			}
		}
	},
	changeIcon: async (_icon) => {
		try {
			let payload = _icon;
			let config = {
				headers: {
					"Content-Type": "image/*",
				},
			};
	
			// If the icon is a base64 string, convert it to a File object
			if (typeof _icon === "string") {
				payload = imageConverter.base64ToFile(_icon);
				if (!payload) {
					throw new Error("Failed to convert image");
				}
			}
	
			const response = await api.put("/settings/profile/icon", payload, config);
			if (response.status === 200) {
				return true;
			}
			return false;
		} catch (error) {
			console.error("Error changing icon:", error);
			throw error;
		}
	}
};
