<template>
  <div class="work-form">
    <div class="form-header">
      <router-link :to="isEdit ? `/works/${id}` : '/'" class="back-link">
        <i class="material-icons">arrow_back</i>
        <span>{{ isEdit ? 'Détail' : 'Liste' }}</span>
      </router-link>
    </div>

    <div class="form-card qtm-card">
      <div class="card-header">
        <h1 class="card-title">
          <i class="material-icons">{{ isEdit ? 'edit' : 'add_circle_outline' }}</i>
          {{ isEdit ? 'Modifier l\'œuvre' : 'Nouvelle œuvre' }}
        </h1>
      </div>

      <form class="card-body" @submit.prevent="submit">
        <div class="form-grid">
          <div class="qtm-form-field">
            <label class="qtm-form-label" for="type">
              Type <span class="required">*</span>
            </label>
            <select id="type" v-model="form.type" class="qtm-form-select" required>
              <option v-for="t in WORK_TYPES" :key="t.value" :value="t.value">{{ t.label }}</option>
            </select>
          </div>

          <div class="qtm-form-field">
            <label class="qtm-form-label" for="title">
              Titre <span class="required">*</span>
            </label>
            <input id="title" v-model="form.title" type="text" class="qtm-form-input" required />
          </div>

          <div class="qtm-form-field field-full">
            <label class="qtm-form-label">
              Auteur(s) <span class="required">*</span>
              <span class="label-hint">un par ligne</span>
            </label>
            <textarea v-model="authorsText" rows="3" class="qtm-form-textarea" placeholder="Un auteur par ligne"></textarea>
          </div>

          <div class="qtm-form-field">
            <label class="qtm-form-label" for="origin">Origine</label>
            <input id="origin" v-model="form.origin" type="text" class="qtm-form-input" />
          </div>

          <div class="qtm-form-field">
            <label class="qtm-form-label" for="availability">Disponibilité</label>
            <input id="availability" v-model="form.availability" type="text" class="qtm-form-input" />
          </div>

          <div class="qtm-form-field field-full">
            <label class="checkbox-label">
              <input v-model="form.seen" type="checkbox" class="checkbox-input" />
              <span class="checkbox-custom">
                <i class="material-icons">{{ form.seen ? 'check_box' : 'check_box_outline_blank' }}</i>
              </span>
              <span>Déjà vu</span>
            </label>
          </div>
        </div>

        <div v-if="error" class="qtm-alert qtm-alert-error form-error">
          <i class="material-icons">error_outline</i>
          <span>{{ error }}</span>
        </div>

        <div class="form-actions">
          <button type="submit" class="qtm-btn qtm-btn-primary">
            <i class="material-icons">{{ isEdit ? 'save' : 'add' }}</i>
            <span>{{ isEdit ? 'Enregistrer' : 'Créer' }}</span>
          </button>
          <router-link :to="isEdit ? `/works/${id}` : '/'" class="qtm-btn qtm-btn-outline">
            Annuler
          </router-link>
        </div>
      </form>
    </div>
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
.form-header {
  margin-bottom: var(--qtm-space-xl);
}
.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--qtm-space-xs);
  color: var(--qtm-text-secondary);
  text-decoration: none;
  font-size: var(--qtm-font-size-base);
  font-weight: 500;
  padding: var(--qtm-space-xs) var(--qtm-space-s);
  border-radius: var(--qtm-radius-m);
  transition: all var(--qtm-transition-fast);
  margin-left: calc(-1 * var(--qtm-space-s));
}
.back-link:hover {
  color: var(--qtm-primary-400);
  background: var(--qtm-primary-100);
  text-decoration: none;
}
.back-link .material-icons {
  font-size: 1.125rem;
}

.form-card {
  max-width: 40rem;
  overflow: hidden;
}
.card-header {
  padding: var(--qtm-space-xl) var(--qtm-space-xxl);
  border-bottom: 1px solid var(--qtm-border-default);
  background: var(--qtm-bluegrey-50);
}
.card-title {
  display: flex;
  align-items: center;
  gap: var(--qtm-space-s);
  margin: 0;
  font-size: var(--qtm-font-size-xl);
  font-weight: 700;
  color: var(--qtm-text-primary);
}
.card-title .material-icons {
  color: var(--qtm-primary-400);
}

.card-body {
  padding: var(--qtm-space-xxl);
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 var(--qtm-space-xl);
}
.field-full {
  grid-column: 1 / -1;
}
.required {
  color: var(--qtm-danger-400);
}
.label-hint {
  font-weight: 400;
  font-size: var(--qtm-font-size-sm);
  color: var(--qtm-bluegrey-400);
  text-transform: none;
  letter-spacing: normal;
  margin-left: var(--qtm-space-xs);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--qtm-space-s);
  cursor: pointer;
  font-size: var(--qtm-font-size-base);
  color: var(--qtm-text-primary);
  user-select: none;
}
.checkbox-input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}
.checkbox-custom .material-icons {
  font-size: 1.5rem;
  color: var(--qtm-bluegrey-400);
  transition: color var(--qtm-transition-fast);
}
.checkbox-input:checked + .checkbox-custom .material-icons {
  color: var(--qtm-primary-400);
}

.form-error {
  margin-bottom: var(--qtm-space-xl);
}
.form-actions {
  display: flex;
  gap: var(--qtm-space-s);
  padding-top: var(--qtm-space-l);
  border-top: 1px solid var(--qtm-border-default);
}

@media (max-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
