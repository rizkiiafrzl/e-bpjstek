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

    console.log('API Service - Making request:', {
      url,
      method: config.method || 'GET',
      headers: config.headers,
      body: config.body
    })

    try {
      const response = await fetch(url, config)
      console.log('API Service - Response received:', {
        status: response.status,
        statusText: response.statusText,
        headers: Object.fromEntries(response.headers.entries())
      })
      
      // gracefully handle empty bodies (e.g., 204 No Content)
      const contentType = response.headers.get('content-type') || ''
      let data
      if (response.status === 204 || contentType.indexOf('application/json') === -1) {
        const text = await response.text()
        data = text ? JSON.parse(text) : null
      } else {
        const text = await response.text()
        data = text ? JSON.parse(text) : null
      }

      if (!response.ok) {
        const message = (data && (data.error || data.message)) || `HTTP error! status: ${response.status}`
        console.error('API Service - Request failed:', {
          status: response.status,
          statusText: response.statusText,
          message,
          data
        })
        throw new Error(message)
      }

      console.log('API Service - Request successful:', data)
      // return parsed json (or null) so callers can decide
      return data
    } catch (error) {
      console.error('API Service - Request error:', {
        message: error.message,
        stack: error.stack,
        url,
        config
      })
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
  async getReportPeriods({ page = 1, pageSize = 10, status = 'all' } = {}) {
    const qs = new URLSearchParams({ page: String(page), pageSize: String(pageSize), status })
    return this.request(`/report-periods?${qs.toString()}`, { method: 'GET' })
  }

  async createReportPeriod(payload = {}) {
    return this.request('/report-periods', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  }

  async deleteReportPeriod(id) {
    return this.request(`/report-periods/${id}`, { method: 'DELETE' })
  }

  async calculateReportPeriod(id) {
    return this.request(`/report-periods/${id}/calculate`, { method: 'POST' })
  }

  async finalizeReportPeriod(id) {
    return this.request(`/report-periods/${id}/finalize`, { method: 'POST' })
  }

  async printReportPeriod(id) {
    return this.request(`/report-periods/${id}/print`, { method: 'GET' })
  }

  // Dashboard summary
  async getReportSummary() {
    return this.request('/report-periods-summary', { method: 'GET' })
  }

  // Upload Workers (mass upload)
  async uploadWorkers(file) {
    const form = new FormData()
    form.append('file', file)
    // Don't set Content-Type explicitly; let browser set multipart boundary
    return fetch(`${this.baseURL}/workers/upload`, {
      method: 'POST',
      headers: this.token ? { Authorization: `Bearer ${this.token}` } : {},
      body: form,
    }).then(async (res) => {
      const text = await res.text()
      const data = text ? JSON.parse(text) : null
      if (!res.ok) {
        const message = (data && (data.error || data.message)) || `HTTP error! status: ${res.status}`
        throw new Error(message)
      }
      return data
    })
  }

  // Upload TK (Tenaga Kerja)
  async uploadTK(file) {
    const form = new FormData()
    form.append('file', file)
    // Don't set Content-Type explicitly; let browser set multipart boundary
    return fetch(`${this.baseURL}/workers/upload-tk`, {
      method: 'POST',
      headers: this.token ? { Authorization: `Bearer ${this.token}` } : {},
      body: form,
    }).then(async (res) => {
      const text = await res.text()
      const data = text ? JSON.parse(text) : null
      if (!res.ok) {
        const message = (data && (data.error || data.message)) || `HTTP error! status: ${res.status}`
        throw new Error(message)
      }
      return data
    })
  }

  // Upload TK NA (Tenaga Kerja Nonaktif)
  async uploadTKNA(file) {
    const form = new FormData()
    form.append('file', file)
    const url = `${this.baseURL}/workers/upload-tk-na`
    const headers = this.token ? { Authorization: `Bearer ${this.token}` } : {}
    return fetch(url, {
      method: 'POST',
      headers,
      body: form,
    }).then(async (res) => {
      const text = await res.text()
      let data = null
      try { data = text ? JSON.parse(text) : null } catch { data = { message: text } }
      if (!res.ok) {
        const message = (data && (data.error || data.message)) || `HTTP error! status: ${res.status}`
        throw new Error(message)
      }
      return data
    })
  }

  // Upload Upah (Wage Data)
  async uploadUpah(file) {
    const form = new FormData()
    form.append('file', file)
    // Don't set Content-Type explicitly; let browser set multipart boundary
    return fetch(`${this.baseURL}/workers/upload-upah`, {
      method: 'POST',
      headers: this.token ? { Authorization: `Bearer ${this.token}` } : {},
      body: form,
    }).then(async (res) => {
      const text = await res.text()
      const data = text ? JSON.parse(text) : null
      if (!res.ok) {
        const message = (data && (data.error || data.message)) || `HTTP error! status: ${res.status}`
        throw new Error(message)
      }
      return data
    })
  }

  // Koreksi TK (Correction of Worker Data)
  async koreksiTK(file) {
    const form = new FormData()
    form.append('file', file)
    // Don't set Content-Type explicitly; let browser set multipart boundary
    return fetch(`${this.baseURL}/koreksi-data-tk`, {
      method: 'POST',
      headers: this.token ? { Authorization: `Bearer ${this.token}` } : {},
      body: form,
    }).then(async (res) => {
      const text = await res.text()
      const data = text ? JSON.parse(text) : null
      if (!res.ok) {
        const message = (data && (data.error || data.message)) || `HTTP error! status: ${res.status}`
        throw new Error(message)
      }
      return data
    })
  }

  async getUploadHistory() {
    return this.request('/workers/upload-history', { method: 'GET' })
  }

  async downloadErrorFile(id) {
    const url = `${this.baseURL}/workers/upload-history/${id}/download`
    const config = {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${this.token}`,
      },
    }
    try {
      const response = await fetch(url, config)
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`)
      }
      return response.blob() // Return blob for file download
    } catch (error) {
      console.error('API request failed:', error)
      throw error
    }
  }

  async deleteUploadHistory(id) {
    console.log('API Service - deleteUploadHistory called with ID:', id)
    console.log('API Service - baseURL:', this.baseURL)
    console.log('API Service - token:', this.token ? 'Present' : 'Missing')
    
    const result = await this.request(`/workers/upload-history/${id}`, {
      method: 'DELETE',
    })
    
    console.log('API Service - deleteUploadHistory result:', result)
    return result
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
