export const imageConverter = {
    base64ToFile: (base64String, filename = 'forwarded-image.jpg', type = 'image/jpeg') => {
      try {
        const byteCharacters = atob(base64String)
        const byteNumbers = new Array(byteCharacters.length)
        
        for (let i = 0; i < byteCharacters.length; i++) {
          byteNumbers[i] = byteCharacters.charCodeAt(i)
        }
        
        const byteArray = new Uint8Array(byteNumbers)
        const blob = new Blob([byteArray], { type })
        
        return new File([blob], filename, { type })
      } catch (error) {
        console.error('Error converting base64 to File:', error)
        return null
      }
    }
  }