const API_BASE = 'http://localhost:8080/api/v1';

function toggleChannelFields() {
    const channel = document.getElementById('channel').value;
    const telegramFields = document.getElementById('telegramFields');
    const emailFields = document.getElementById('emailFields');

    if (channel === 'telegram') {
        telegramFields.style.display = 'block';
        emailFields.style.display = 'none';

        document.getElementById('telegramRecipient').required = true;
        document.getElementById('emailRecipient').required = false;
    } else if (channel === 'email') {
        telegramFields.style.display = 'none';
        emailFields.style.display = 'block';

        document.getElementById('telegramRecipient').required = false;
        document.getElementById('emailRecipient').required = true;
    }
}

function getRecipientId() {
    const channel = document.getElementById('channel').value;
    if (channel === 'telegram') {
        return document.getElementById('telegramRecipient').value;
    } else if (channel === 'email') {
        return document.getElementById('emailRecipient').value;
    }
    return '';
}

function validateEmailFields() {
    const channel = document.getElementById('channel').value;
    if (channel !== 'email') return true;

    const requiredFields = [
        'emailRecipient',
        'emailFromEmail',
        'emailUsername',
        'emailPassword'
    ];

    for (const fieldId of requiredFields) {
        const field = document.getElementById(fieldId);
        if (!field.value.trim()) {
            alert(`Please fill in the ${field.placeholder || fieldId} field`);
            field.focus();
            return false;
        }
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    const recipientEmail = document.getElementById('emailRecipient').value;
    const fromEmail = document.getElementById('emailFromEmail').value;

    if (!emailRegex.test(recipientEmail)) {
        alert('Please enter a valid recipient email address');
        document.getElementById('emailRecipient').focus();
        return false;
    }

    if (!emailRegex.test(fromEmail)) {
        alert('Please enter a valid sender email address');
        document.getElementById('emailFromEmail').focus();
        return false;
    }

    return true;
}

document.getElementById('createForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    if (!validateEmailFields()) {
        return;
    }

    const notification = {
        payload: document.getElementById('payload').value,
        notification_date: new Date(document.getElementById('notificationDate').value).toISOString(),
        recipient_id: getRecipientId(),
        channel: document.getElementById('channel').value
    };

    if (notification.channel === 'email') {
        notification.email_config = {
            subject: document.getElementById('emailSubject').value,
            from_name: document.getElementById('emailFromName').value,
            from_email: document.getElementById('emailFromEmail').value,
            smtp_host: document.getElementById('emailSMTPHost').value,
            smtp_port: parseInt(document.getElementById('emailSMTPPort').value),
            username: document.getElementById('emailUsername').value,
            password: document.getElementById('emailPassword').value
        };
    }

    try {
        console.log('Sending request to:', `${API_BASE}/notify`);
        console.log('Request data:', notification);

        const response = await fetch(`${API_BASE}/notify`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(notification)
        });

        console.log('Response status:', response.status);
        console.log('Response headers:', response.headers);

        if (response.ok) {
            const result = await response.json();
            console.log('Response data:', result);
            alert(`Notification created! ID: ${result.result.id}`);
            document.getElementById('createForm').reset();
            toggleChannelFields();
        } else {
            const error = await response.json();
            console.error('Error response:', error);
            alert(`Error creating notification: ${error.error}`);
        }
    } catch (error) {
        console.error('Network error:', error);
        alert(`Network error: ${error.message}`);
    }
});

document.addEventListener('DOMContentLoaded', function () {
    toggleChannelFields();
});

async function findNotification() {
    const id = document.getElementById('searchId').value.trim();
    if (!id) {
        alert('Please enter a notification ID');
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/notify/${id}`);
        const result = document.getElementById('searchResult');

        if (response.ok) {
            const data = await response.json();
            const notification = data.result;
            const status = notification.status;

            result.innerHTML = `
                <div class="notification-result">
                    <div class="status-display">
                        <strong>Status:</strong> 
                        <span class="status-${status}">${status.toUpperCase()}</span>
                    </div>
                    <div class="notification-actions">
                        <button onclick="cancelNotification('${id}')" 
                                class="cancel-btn" 
                                ${status === 'cancelled' || status === 'sent' ? 'disabled' : ''}>
                            Cancel Notification
                        </button>
                    </div>
                </div>
            `;
        } else {
            const error = await response.json();
            result.innerHTML = `<div class="error">${error.error}</div>`;
        }
    } catch (error) {
        document.getElementById('searchResult').innerHTML = '<div class="error">Error fetching notification</div>';
    }
}

async function cancelNotification(id) {
    if (!confirm('Are you sure you want to cancel this notification?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/notify/${id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            alert('Notification cancelled successfully!');
            findNotification();
        } else {
            const error = await response.json();
            alert(`Error cancelling notification: ${error.error}`);
        }
    } catch (error) {
        alert('Error cancelling notification');
    }
}