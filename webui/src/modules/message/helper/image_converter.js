export const imageConverter = {
	base64ToFile: (
		base64String,
		filename = "forwarded-image.jpg",
		type = "image/jpeg",
	) => {
		try {
			const byteCharacters = atob(base64String);
			const byteNumbers = new Array(byteCharacters.length);

			for (let i = 0; i < byteCharacters.length; i++) {
				byteNumbers[i] = byteCharacters.charCodeAt(i);
			}

			const byteArray = new Uint8Array(byteNumbers);
			const blob = new Blob([byteArray], { type });

			return new File([blob], filename, { type });
		} catch (error) {
			console.error("Error converting base64 to File:", error);
			return null;
		}
	},
	fileToBase64: (file) => {
		return new Promise((resolve, reject) => {
		  const reader = new FileReader();
		  
		  reader.onload = () => {
			try {
			  // Remove the "data:image/jpeg;base64," part from the string
			  const base64String = reader.result.split(',')[1];
			  resolve(base64String);
			} catch (error) {
			  reject(error);
			}
		  };
	
		  reader.onerror = (error) => {
			reject(error);
		  };
	
		  reader.readAsDataURL(file);
		});
	  }
};
