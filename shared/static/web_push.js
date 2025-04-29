async function requestNotificationPermission() {
    const permission = await Notification.requestPermission();

    if (permission === 'granted') {
        console.log('Notification permission granted.');
    } else if (permission === 'denied') {
        console.log('Notification permission denied.');
    } else {
        console.log('Notification permission dismissed.');
    }
}

await requestNotificationPermission();