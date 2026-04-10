<template>
  <div class="work-form page-shell">
    <header class="page-header">
      <div>
        <h1 class="page-title">{{ isEdit ? 'Modifier l\'œuvre' : 'Nouvelle œuvre' }}</h1>
        <p class="page-subtitle">
          {{ isEdit ? 'Mettez à jour les informations de l\'œuvre.' : 'Renseignez les champs pour enrichir votre catalogue.' }}
        </p>
      </div>
    </header>

    <form @submit.prevent="submit">
      <div class="form-grid">
        <div class="field">
          <label for="type">Type *</label>
          <select id="type" v-model="form.type" required class="control-select">
            <option v-for="t in WORK_TYPES" :key="t.value" :value="t.value">{{ t.label }}</option>
          </select>
        </div>
        <div class="field">
          <label for="title">Titre *</label>
          <input id="title" v-model="form.title" type="text" required class="control-input" />
        </div>
        <div class="field field-wide">
          <label>Auteur(s) * (un par ligne)</label>
          <textarea
            v-model="authorsText"
            rows="4"
            placeholder="Un auteur par ligne"
            class="control-textarea"
          ></textarea>
        </div>
        <div class="field">
          <label for="origin">Origine</label>
          <input id="origin" v-model="form.origin" type="text" class="control-input" />
        </div>
        <div class="field">
          <label for="availability">Disponibilité</label>
          <input id="availability" v-model="form.availability" type="text" class="control-input" />
        </div>
        <div class="field checkbox field-wide">
          <label>
            <input v-model="form.seen" type="checkbox" />
            Déjà vu
          </label>
        </div>
      </div>

      <p v-if="error" class="inline-error">{{ error }}</p>
      <div class="actions">
        <button type="submit" class="action-btn primary">{{ isEdit ? 'Enregistrer' : 'Créer' }}</button>
        <router-link :to="isEdit ? `/works/${id}` : '/'" class="action-btn">Annuler</router-link>
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
.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem 1.2rem;
}

.field {
  margin-bottom: 0.1rem;
}

.field-wide {
  grid-column: 1 / -1;
}

.field label {
  display: block;
  margin-bottom: 0.42rem;
  color: var(--q-text-muted);
  font-size: 0.76rem;
  letter-spacing: 0.07em;
  text-transform: uppercase;
}

.field.checkbox label {
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
  color: var(--q-text-secondary);
  font-size: 0.9rem;
  text-transform: none;
  letter-spacing: 0;
}

.field.checkbox input {
  width: 1rem;
  height: 1rem;
  accent-color: var(--q-accent-strong);
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.6rem;
  margin-top: 1.5rem;
}
</style>
