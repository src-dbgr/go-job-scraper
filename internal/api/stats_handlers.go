package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (a *API) getTopJobCategories(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetTopJobCategories()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top job categories")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getAvgExperienceByCategory(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetAvgExperienceByCategory()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get average experience by category")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getRemoteVsOnsite(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetRemoteVsOnsite()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get remote vs onsite data")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getTopSkills(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetTopSkills()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top skills")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getTopOptionalSkills(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetTopOptionalSkills()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top optional skills")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getBenefitsByCompanySize(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetBenefitsByCompanySize()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get benefits by company size")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getAvgSalaryByEducation(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetAvgSalaryByEducation()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get average salary by education")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getJobPostingsTrend(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetJobPostingsTrend()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job postings trend")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getLanguagesByLocation(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetLanguagesByLocation()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get languages by location")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getEmploymentTypes(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetEmploymentTypes()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get employment types")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getRemoteWorkByCategory(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetRemoteWorkByCategory()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get remote work by category")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getTechnologyTrends(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetTechnologyTrends()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get technology trends")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getJobRequirementsByLocation(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetJobRequirementsByLocation()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job requirements by location")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getRemoteVsOnsiteByIndustry(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetRemoteVsOnsiteByIndustry()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get remote vs onsite by industry")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getJobCategoriesByCompanySize(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetJobCategoriesByCompanySize()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job categories by company size")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getSkillsByExperienceLevel(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetSkillsByExperienceLevel()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get skills by experience level")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getCompaniesBySize(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetCompaniesBySize()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get companies by size")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getCompaniesBySizeAndType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sizeType := vars["sizeType"]
	result, err := a.jobStatsService.GetCompaniesBySizeAndType(sizeType)
	if err != nil {
		log.Error().Err(err).Str("sizeType", sizeType).Msg("Failed to get companies by size and type")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getCompanySizeDistribution(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetCompanySizeDistribution()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get company size distribution")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getJobPostingsPerDay(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetJobPostingsPerDay()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job postings per day")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getJobPostingsPerMonth(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetJobPostingsPerMonth()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job postings per month")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getJobPostingsPerCompany(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetJobPostingsPerCompany()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job postings per company")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getMustSkillFrequencyPerDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skill := vars["skill"]
	result, err := a.jobStatsService.GetMustSkillFrequencyPerDay(skill)
	if err != nil {
		log.Error().Err(err).Str("skill", skill).Msg("Failed to get must skill frequency per day")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getOptionalSkillFrequencyPerDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skill := vars["skill"]
	result, err := a.jobStatsService.GetOptionalSkillFrequencyPerDay(skill)
	if err != nil {
		log.Error().Err(err).Str("skill", skill).Msg("Failed to get optional skill frequency per day")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}

func (a *API) getJobCategoryCounts(w http.ResponseWriter, r *http.Request) {
	result, err := a.jobStatsService.GetJobCategoryCounts()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job category counts")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
}
