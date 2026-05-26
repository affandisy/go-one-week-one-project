import { create } from 'zustand';

export const useBillingStore = create((set, get) => ({
  // Data pemilik yang sedang dilayani
  currentOwner: null,
  
  // Daftar hewan peliharaan milik owner tersebut
  ownerPets: [], 
  
  // Keranjang tagihan: array of { id, pet_id, service_id, price, pet_name, service_name }
  cartItems: [],

  // Aksi untuk mengatur sesi transaksi
  startTransaction: (owner, pets) => set({ 
    currentOwner: owner, 
    ownerPets: pets, 
    cartItems: [] 
  }),

  clearTransaction: () => set({ 
    currentOwner: null, 
    ownerPets: [], 
    cartItems: [] 
  }),

  // FR-004: Menambahkan layanan spesifik untuk hewan spesifik
  addItem: (pet, service) => set((state) => ({
    cartItems: [
      ...state.cartItems,
      {
        id: crypto.randomUUID(), // ID sementara untuk UI
        pet_id: pet.id,
        pet_name: pet.name,
        service_id: service.id,
        service_name: service.name,
        price: service.base_price
      }
    ]
  })),

  // Menghapus item dari keranjang
  removeItem: (itemId) => set((state) => ({
    cartItems: state.cartItems.filter(item => item.id !== itemId)
  })),

  // Menghitung total keseluruhan
  getTotal: () => {
    return get().cartItems.reduce((total, item) => total + item.price, 0);
  },

  // Mengelompokkan item berdasarkan hewan untuk tampilan struk/UI yang rapi
  getItemsGroupedByPet: () => {
    const items = get().cartItems;
    const grouped = {};
    
    items.forEach(item => {
      if (!grouped[item.pet_id]) {
        grouped[item.pet_id] = { pet_name: item.pet_name, services: [], subtotal: 0 };
      }
      grouped[item.pet_id].services.push(item);
      grouped[item.pet_id].subtotal += item.price;
    });
    
    return grouped;
  }
}));