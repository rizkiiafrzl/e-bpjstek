// JavaScript untuk debug frontend upload massal
// Jalankan di console browser (F12)

console.log("🐛 Debugging Frontend Upload Massal Issue...");

// 1. Cek data yang diambil dari API
async function checkWorkersData() {
    try {
        const response = await fetch('/api/v1/workers', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
                'Content-Type': 'application/json'
            }
        });
        
        if (!response.ok) {
            console.error('❌ API Error:', response.status, response.statusText);
            return;
        }
        
        const workers = await response.json();
        console.log('📊 Total workers from API:', workers.length);
        console.log('📋 Workers data:', workers);
        
        if (workers.length === 0) {
            console.log('❌ TIDAK ADA DATA WORKERS!');
            console.log('   Kemungkinan:');
            console.log('   - Upload gagal');
            console.log('   - Data tidak tersimpan');
            console.log('   - User ID mismatch');
        } else {
            console.log('✅ Ada data workers');
            console.log('❓ Jika tidak tampil di tabel, kemungkinan:');
            console.log('   - Filter aktif');
            console.log('   - Pagination issue');
            console.log('   - Component tidak re-render');
        }
        
        return workers;
    } catch (error) {
        console.error('❌ Error fetching workers:', error);
    }
}

// 2. Cek upload history
async function checkUploadHistory() {
    try {
        const response = await fetch('/api/v1/workers/upload-history', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
                'Content-Type': 'application/json'
            }
        });
        
        if (!response.ok) {
            console.error('❌ API Error:', response.status, response.statusText);
            return;
        }
        
        const history = await response.json();
        console.log('📊 Upload history:', history);
        
        if (history.length === 0) {
            console.log('❌ TIDAK ADA UPLOAD HISTORY!');
            console.log('   Kemungkinan belum pernah upload atau upload gagal');
        } else {
            console.log('✅ Ada upload history');
            const latest = history[0];
            console.log('📋 Upload terbaru:', latest);
            console.log(`   File: ${latest.fileName}`);
            console.log(`   Type: ${latest.type}`);
            console.log(`   Valid: ${latest.totalValid}`);
            console.log(`   Invalid: ${latest.totalInvalid}`);
            console.log(`   Status: ${latest.validationStatus}`);
        }
        
        return history;
    } catch (error) {
        console.error('❌ Error fetching upload history:', error);
    }
}

// 3. Cek filter dan search
function checkFilters() {
    console.log('🔍 Checking filters and search...');
    
    // Cek apakah ada filter aktif
    const filterSelects = document.querySelectorAll('select');
    filterSelects.forEach((select, index) => {
        if (select.value && select.value !== '') {
            console.log(`⚠️ Filter ${index} aktif:`, select.value);
        }
    });
    
    // Cek search input
    const searchInputs = document.querySelectorAll('input[type="text"]');
    searchInputs.forEach((input, index) => {
        if (input.value && input.value !== '') {
            console.log(`⚠️ Search ${index} aktif:`, input.value);
        }
    });
    
    // Cek pagination
    const pagination = document.querySelector('.pagination, [class*="pagination"]');
    if (pagination) {
        console.log('📄 Pagination found:', pagination);
    }
}

// 4. Cek component state
function checkComponentState() {
    console.log('🔍 Checking component state...');
    
    // Cek apakah ada loading state
    const loadingElements = document.querySelectorAll('[class*="loading"], [class*="spinner"]');
    if (loadingElements.length > 0) {
        console.log('⏳ Loading state detected:', loadingElements);
    }
    
    // Cek error messages
    const errorElements = document.querySelectorAll('[class*="error"], [class*="alert"]');
    if (errorElements.length > 0) {
        console.log('❌ Error messages found:', errorElements);
    }
    
    // Cek table rows
    const tableRows = document.querySelectorAll('table tbody tr');
    console.log('📊 Table rows found:', tableRows.length);
    
    if (tableRows.length === 0) {
        console.log('❌ TIDAK ADA ROWS DI TABEL!');
        console.log('   Kemungkinan:');
        console.log('   - Data kosong');
        console.log('   - Filter menyaring semua data');
        console.log('   - Component tidak render');
    }
}

// 5. Test refresh data
async function testRefresh() {
    console.log('🔄 Testing data refresh...');
    
    // Simulate refresh button click
    const refreshButton = document.querySelector('button[class*="refresh"], button:contains("REFRESH")');
    if (refreshButton) {
        console.log('🔄 Clicking refresh button...');
        refreshButton.click();
        
        // Wait a bit and check data again
        setTimeout(async () => {
            console.log('🔍 Checking data after refresh...');
            await checkWorkersData();
        }, 2000);
    } else {
        console.log('⚠️ Refresh button not found');
    }
}

// 6. Main debug function
async function debugUploadMassal() {
    console.log('🚀 Starting comprehensive debug...');
    
    // Check API data
    const workers = await checkWorkersData();
    
    // Check upload history
    const history = await checkUploadHistory();
    
    // Check frontend filters
    checkFilters();
    
    // Check component state
    checkComponentState();
    
    // Test refresh
    await testRefresh();
    
    console.log('✅ Debug completed!');
    console.log('📋 Summary:');
    console.log(`   - Workers data: ${workers ? workers.length : 'Error'}`);
    console.log(`   - Upload history: ${history ? history.length : 'Error'}`);
    console.log('   - Check console above for detailed issues');
}

// 7. Quick fixes
function quickFixes() {
    console.log('🛠️ Quick fixes to try:');
    console.log('1. Refresh halaman (F5)');
    console.log('2. Clear filter/search');
    console.log('3. Check network tab untuk API calls');
    console.log('4. Restart backend');
    console.log('5. Cek log backend untuk error');
}

// Run debug
debugUploadMassal();

// Export functions for manual use
window.debugUploadMassal = debugUploadMassal;
window.checkWorkersData = checkWorkersData;
window.checkUploadHistory = checkUploadHistory;
window.quickFixes = quickFixes;

console.log('💡 Functions available:');
console.log('   - debugUploadMassal()');
console.log('   - checkWorkersData()');
console.log('   - checkUploadHistory()');
console.log('   - quickFixes()');







