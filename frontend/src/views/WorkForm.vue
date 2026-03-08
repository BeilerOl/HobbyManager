<template>
  <div class="work-form">
    <h1>{{ isEdit ? 'Modifier l\'œuvre' : 'Nouvelle œuvre' }}</h1>
    <form @submit.prevent="submit">
      <div class="field">
        <label for="type">Type *</label>
        <select id="type" v-model="form.type" required>
          <option v-for="t in WORK_TYPES" :key="t.value" :value="t.value">{{ t.label }}</option>
        </select>
      </div>
      <div class="field">
        <label for="title">Titre *</label>
        <input id="title" v-model="form.title" type="text" required />
      </div>
      <div class="field">
        <label>Auteur(s) * (un par ligne)</label>
        <textarea v-model="authorsText" rows="3" placeholder="Un auteur par ligne"></textarea>
      </div>
      <div class="field">
        <label for="origin">Origine</label>
        <input id="origin" v-model="form.origin" type="text" />
      </div>
      <div class="field">
        <label for="availability">Disponibilité</label>
        <input id="availability" v-model="form.availability" type="text" />
      </div>
      <div class="field checkbox">
        <label>
          <input v-model="form.seen" type="checkbox" />
          Déjà vu
        </label>
      </div>
      <p v-if="error" class="error">{{ error }}</p>
      <div class="actions">
        <button type="submit" class="btn primary">{{ isEdit ? 'Enregistrer' : 'Créer' }}</button>
        <router-link :to="isEdit ? `/works/${id}` : '/'" class="btn">Annuler</router-link>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getWork, createWork, updateWork, WORK_TYPES } from '../api/works'

const props = defineProps({ isEdit: Boolean })
const route = useRoute()
const router = useRouter()

const form = ref({
  type: 'roman',
  title: '',
  authors: [],
  origin: '',
  availability: '',
  seen: false,
})
const authorsText = ref('')
const error = ref('')
const loading = ref(false)

const id = computed(() => Number(route.params.id))

// Sync authors array with textarea (one author per line)
watch(authorsText, (t) => {
  form.value.authors = t
    .split('\n')
    .map((s) => s.trim())
    .filter(Boolean)
})

async function load() {
  if (!props.isEdit) return
  loading.value = true
  error.value = ''
  try {
    const w = await getWork(id.value)
    if (w) {
      form.value = {
        type: w.type,
        title: w.title,
        authors: w.authors || [],
        origin: w.origin || '',
        availability: w.availability || '',
        seen: w.seen ?? false,
      }
      authorsText.value = (w.authors || []).join('\n')
    } else {
      error.value = 'Œuvre introuvable.'
    }
  } catch (e) {
    error.value = e.message || 'Erreur'
  } finally {
    loading.value = false
  }
}

async function submit() {
  error.value = ''
  const payload = {
    type: form.value.type,
    title: form.value.title.trim(),
    authors: form.value.authors.length ? form.value.authors : [],
    origin: form.value.origin.trim(),
    availability: form.value.availability.trim(),
    seen: form.value.seen,
  }
  try {
    if (props.isEdit) {
      await updateWork(id.value, payload)
      router.push(`/works/${id.value}`)
    } else {
      const created = await createWork(payload)
      router.push(`/works/${created.id}`)
    }
  } catch (e) {
    error.value = e.message || 'Erreur'
  }
}

onMounted(load)
</script>

<style scoped>
.work-form h1 {
  margin-top: 0;
  font-size: 1.5rem;
}
.field {
  margin-bottom: 1rem;
}
.field label {
  display: block;
  margin-bottom: 0.25rem;
  color: #a1a1aa;
  font-size: 0.9rem;
}
.field input[type="text"],
.field select,
.field textarea {
  width: 100%;
  max-width: 28rem;
  padding: 0.5rem 0.75rem;
  background: #25262b;
  border: 1px solid #3f3f46;
  border-radius: 6px;
  color: #e4e4e7;
}
.field.checkbox label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: inherit;
}
.field.checkbox input {
  width: auto;
}
.actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 1.5rem;
}
.btn {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: 1px solid #3f3f46;
  background: #25262b;
  color: #e4e4e7;
  cursor: pointer;
  text-decoration: none;
  font-size: 0.9rem;
}
.btn:hover {
  background: #3f3f46;
}
.btn.primary {
  border-color: #3b82f6;
  background: #2563eb;
  color: #fff;
}
.btn.primary:hover {
  background: #1d4ed8;
}
.error {
  color: #f87171;
  margin-bottom: 1rem;
}
</style>
