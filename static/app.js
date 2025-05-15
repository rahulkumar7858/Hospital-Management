let token = "";
let userRole = "";

function showForm(formType) {
    document.getElementById('login-form').style.display = formType === 'login' ? 'block' : 'none';
    document.getElementById('signup-form').style.display = formType === 'signup' ? 'block' : 'none';
    document.getElementById('login-tab').classList.toggle('active', formType === 'login');
    document.getElementById('signup-tab').classList.toggle('active', formType === 'signup');
    document.getElementById('login-error').innerText = "";
    document.getElementById('signup-error').innerText = "";
    document.getElementById('signup-success').innerText = "";
}

async function login(e) {
    e.preventDefault();
    let username = document.getElementById('login-username').value;
    let password = document.getElementById('login-password').value;
    let role = document.getElementById('login-role').value;
    let res = await fetch("/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password, role })
    });
    let data = await res.json();
    if (data.token) {
        token = data.token;
        userRole = role;
        document.getElementById("auth-container").style.display = "none";
        document.getElementById("dashboard").style.display = "block";
        showDashboard();
    } else {
        document.getElementById("login-error").innerText = data.error || "Login failed";
    }
}

async function signup(e) {
    e.preventDefault();
    let username = document.getElementById('signup-username').value;
    let password = document.getElementById('signup-password').value;
    let role = document.getElementById('signup-role').value;
    let res = await fetch("/signup", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password, role })
    });
    let data = await res.json();
    if (res.ok) {
        document.getElementById("signup-success").innerText = data.message || "Signup successful!";
        document.getElementById("signup-error").innerText = "";
        setTimeout(() => showForm('login'), 1200);
    } else {
        document.getElementById("signup-error").innerText = data.error || "Signup failed";
        document.getElementById("signup-success").innerText = "";
    }
}

function showDashboard() {
    if (userRole === "receptionist") {
        document.getElementById("reception-portal").style.display = "block";
        document.getElementById("doctor-portal").style.display = "none";
        loadReceptionistPatients();
    } else {
        document.getElementById("reception-portal").style.display = "none";
        document.getElementById("doctor-portal").style.display = "block";
        loadPatients();
    }
}

async function createPatient(e) {
    e.preventDefault();
    let name = document.getElementById("pname").value;
    let age = Number(document.getElementById("page").value);
    let gender = document.getElementById("pgender").value;
    let diagnosis = "";
    let res = await fetch("/reception/patients", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify({ name, age, gender, diagnosis })
    });
    let data = await res.json();
    document.getElementById("createResult").innerText = res.ok ? `Patient created: ${JSON.stringify(data)}` : (data.error || "Error");
    if (res.ok) {
        document.getElementById("createPatientForm").reset();
        loadReceptionistPatients();
    }
}

async function loadReceptionistPatients() {
    let res = await fetch("/reception/patients", {
        headers: { "Authorization": "Bearer " + token }
    });
    let patients = await res.json();
    let html = '<div class="patient-list">';
    patients.forEach(p => {
        html += `
            <div class="patient-card">
                <div class="patient-header">
                    <span class="patient-id">#${p.id}</span>
                    <h3 class="patient-name">${p.name}</h3>
                </div>
                <div class="patient-details">
                    <p>Age: ${p.age}</p>
                    <p>Gender: ${p.gender}</p>
                    <p class="diagnosis">Diagnosis: ${p.diagnosis || "None"}</p>
                    <button class="edit-btn" onclick="showEditPatient(${p.id}, '${p.name}', ${p.age}, '${p.gender}', '${p.diagnosis || ""}')">Edit</button>
                    <button class="delete-btn" onclick="deletePatient(${p.id})">Delete</button>
                </div>
            </div>
        `;
    });
    html += '</div>';
    document.getElementById("receptionPatientsList").innerHTML = html;
}

function showEditPatient(id, name, age, gender, diagnosis) {
    document.getElementById("editPatientId").value = id;
    document.getElementById("edit-pname").value = name;
    document.getElementById("edit-page").value = age;
    document.getElementById("edit-pgender").value = gender;
    document.getElementById("edit-pdiagnosis").value = diagnosis;
    document.getElementById("editPatientModal").style.display = "flex";
}

async function updatePatient(e) {
    e.preventDefault();
    let id = document.getElementById("editPatientId").value;
    let name = document.getElementById("edit-pname").value;
    let age = Number(document.getElementById("edit-page").value);
    let gender = document.getElementById("edit-pgender").value;
    let diagnosis = document.getElementById("edit-pdiagnosis").value;
    let res = await fetch(`/reception/patients/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify({ name, age, gender, diagnosis })
    });
    let data = await res.json();
    document.getElementById("editResult").innerText = res.ok ? `Patient updated: ${JSON.stringify(data)}` : (data.error || "Error");
    if (res.ok) {
        document.getElementById("editPatientModal").style.display = "none";
        loadReceptionistPatients();
    }
}

async function deletePatient(id) {
    let res = await fetch(`/reception/patients/${id}`, {
        method: "DELETE",
        headers: { "Authorization": "Bearer " + token }
    });
    let data = await res.json();
    loadReceptionistPatients();
}

async function loadPatients() {
    let res = await fetch("/doctor/patients", {
        headers: { "Authorization": "Bearer " + token }
    });
    let patients = await res.json();
    let html = '<div class="patient-list">';
    patients.forEach(p => {
        html += `
            <div class="patient-card">
                <div class="patient-header">
                    <span class="patient-id">#${p.id}</span>
                    <h3 class="patient-name">${p.name}</h3>
                </div>
                <div class="patient-details">
                    <p>Age: ${p.age}</p>
                    <p>Gender: ${p.gender}</p>
                    <p class="diagnosis">Diagnosis: ${p.diagnosis || "None"}</p>
                </div>
            </div>
        `;
    });
    html += '</div>';
    document.getElementById("patientsList").innerHTML = html;
}

function closeEditModal() {
    document.getElementById("editPatientModal").style.display = "none";
}

function logout() {
    token = "";
    userRole = "";
    document.getElementById("dashboard").style.display = "none";
    document.getElementById("auth-container").style.display = "flex";
    showForm('login');
}

showForm('login');
