// API Service for Farm Management
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

class ApiService {
  constructor() {
    this.baseURL = API_BASE_URL
    this.token = localStorage.getItem('token')
  }

  // Set authentication token
  setToken(token) {
    this.token = token
    localStorage.setItem('token', token)
  }

  // Clear authentication token
  clearToken() {
    this.token = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('isLoggedIn')
  }

  // Get headers for API requests
  getHeaders() {
    const headers = {
      'Content-Type': 'application/json',
    }

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    return headers
  }

  // Make HTTP request
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`
    const config = {
      headers: this.getHeaders(),
      ...options,
    }

    try {
      const response = await fetch(url, config)
      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || `HTTP error! status: ${response.status}`)
      }

      return data
    } catch (error) {
      console.error('API request failed:', error)
      throw error
    }
  }

  // Authentication endpoints
  async register(userData) {
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    })
  }

  async login(credentials) {
    const response = await this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    })

    // Store token and user data
    if (response.token) {
      this.setToken(response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      localStorage.setItem('isLoggedIn', 'true')
      localStorage.setItem('loginTime', Date.now().toString())
    }

    return response
  }

  async logout() {
    try {
      await this.request('/logout', {
        method: 'POST',
      })
    } finally {
      this.clearToken()
    }
  }

  // User profile endpoints
  async getProfile() {
    return this.request('/profile')
  }

  async updateProfile(profileData) {
    return this.request('/profile', {
      method: 'PUT',
      body: JSON.stringify(profileData),
    })
  }

  // Farm management endpoints
  async getFarms() {
    return this.request('/farms')
  }

  async createFarm(farmData) {
    return this.request('/farms', {
      method: 'POST',
      body: JSON.stringify(farmData),
    })
  }

  async getFarm(id) {
    return this.request(`/farms/${id}`)
  }

  async updateFarm(id, farmData) {
    return this.request(`/farms/${id}`, {
      method: 'PUT',
      body: JSON.stringify(farmData),
    })
  }

  async deleteFarm(id) {
    return this.request(`/farms/${id}`, {
      method: 'DELETE',
    })
  }

  // reCAPTCHA validation
  async validateRecaptcha(recaptchaToken) {
    return this.request('/auth/validate-recaptcha', {
      method: 'POST',
      body: JSON.stringify({ recaptchaToken }),
    })
  }

  // Health check
  async healthCheck() {
    return this.request('/health', {
      method: 'GET',
    })
  }

  // Report Periods
  async getReportPeriods() {
    return this.request('/report-periods', {
      method: 'GET',
    })
  }

  async createReportPeriod(payload = {}) {
    return this.request('/report-periods', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  }

  // Workers
  async getWorkers() {
    return this.request('/workers', { method: 'GET' })
  }

  async getWorker(id) {
    return this.request(`/workers/${id}`, { method: 'GET' })
  }

  async createWorker(worker) {
    return this.request('/workers', {
      method: 'POST',
      body: JSON.stringify(worker),
    })
  }

  async updateWorker(id, worker) {
    return this.request(`/workers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(worker),
    })
  }

  async deleteWorker(id) {
    return this.request(`/workers/${id}`, {
      method: 'DELETE',
    })
  }
}

// Create singleton instance
const apiService = new ApiService()

export default apiService
