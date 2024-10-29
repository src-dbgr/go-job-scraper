package services

import (
	"context"
	"errors"
	"math"
	"time"

	"job-scraper/internal/apperrors"
	"job-scraper/internal/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobStatisticsService struct {
	storage storage.Storage
}

func NewJobStatisticsService(storage storage.Storage) *JobStatisticsService {
	return &JobStatisticsService{storage: storage}
}

func (s *JobStatisticsService) GetTopJobCategories() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$jobCategories"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$jobCategories"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		{{Key: "$limit", Value: 10}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetAvgExperienceByCategory() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$jobCategories"}},
		{{Key: "$match", Value: bson.D{{Key: "yearsOfExperience", Value: bson.D{{Key: "$gt", Value: 0}}}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$jobCategories"},
			{Key: "avgExperience", Value: bson.D{{Key: "$avg", Value: "$yearsOfExperience"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "avgExperience", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetRemoteVsOnsite() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$remote"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "workType", Value: bson.D{{Key: "$cond", Value: bson.D{
				{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$_id", true}}}},
				{Key: "then", Value: "Remote"},
				{Key: "else", Value: "On-site"},
			}}}},
			{Key: "count", Value: 1},
		}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetTopSkills() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$mustSkills"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$mustSkills"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		{{Key: "$limit", Value: 100}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetTopOptionalSkills() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$optionalSkills"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$optionalSkills"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		{{Key: "$limit", Value: 100}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetBenefitsByCompanySize() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$benefits"}},
		{{Key: "$project", Value: bson.D{
			{Key: "companySize", Value: bson.D{{Key: "$switch", Value: bson.D{
				{Key: "branches", Value: bson.A{
					bson.D{{Key: "case", Value: bson.D{{Key: "$lte", Value: bson.A{"$companySize", 50}}}}, {Key: "then", Value: "Small"}},
					bson.D{{Key: "case", Value: bson.D{{Key: "$lte", Value: bson.A{"$companySize", 250}}}}, {Key: "then", Value: "Medium"}},
				}},
				{Key: "default", Value: "Large"},
			}}}},
			{Key: "benefit", Value: "$benefits"},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "companySize", Value: "$companySize"}, {Key: "benefit", Value: "$benefit"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetAvgSalaryByEducation() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "salary", Value: bson.D{
				{Key: "$ne", Value: nil},
				{Key: "$ne", Value: ""},
				{Key: "$gt", Value: 0},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$educationLevel"},
			{Key: "avgSalary", Value: bson.D{{Key: "$avg", Value: "$salary"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "avgSalary", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetJobPostingsTrend() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$project", Value: bson.D{
			{Key: "monthYear", Value: bson.D{{Key: "$dateToString", Value: bson.D{
				{Key: "format", Value: "%Y-%m"},
				{Key: "date", Value: "$postingDate"},
			}}}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$monthYear"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetLanguagesByLocation() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$languages"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "location", Value: "$location"}, {Key: "language", Value: "$languages"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetEmploymentTypes() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$employmentType"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetRemoteWorkByCategory() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$jobCategories"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "category", Value: "$jobCategories"}, {Key: "remote", Value: "$remote"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetTechnologyTrends() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$mustSkills"}},
		{{Key: "$project", Value: bson.D{
			{Key: "monthYear", Value: bson.D{{Key: "$dateToString", Value: bson.D{
				{Key: "format", Value: "%Y-%m"},
				{Key: "date", Value: "$postingDate"},
			}}}},
			{Key: "skill", Value: "$mustSkills"},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "monthYear", Value: "$monthYear"}, {Key: "skill", Value: "$skill"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id.monthYear", Value: 1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetJobRequirementsByLocation() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$mustSkills"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "location", Value: "$location"}, {Key: "skill", Value: "$mustSkills"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetRemoteVsOnsiteByIndustry() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$jobCategories"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "category", Value: "$jobCategories"}, {Key: "remote", Value: "$remote"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "category", Value: "$_id.category"},
			{Key: "workType", Value: bson.D{{Key: "$cond", Value: bson.D{
				{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$_id.remote", true}}}},
				{Key: "then", Value: "Remote"},
				{Key: "else", Value: "On-site"},
			}}}},
			{Key: "count", Value: 1},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetJobCategoriesByCompanySize() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$jobCategories"}},
		{{Key: "$project", Value: bson.D{
			{Key: "companySize", Value: bson.D{{Key: "$switch", Value: bson.D{
				{Key: "branches", Value: bson.A{
					bson.D{{Key: "case", Value: bson.D{{Key: "$lte", Value: bson.A{"$companySize", 50}}}}, {Key: "then", Value: "Small"}},
					bson.D{{Key: "case", Value: bson.D{{Key: "$lte", Value: bson.A{"$companySize", 250}}}}, {Key: "then", Value: "Medium"}},
				}},
				{Key: "default", Value: "Large"},
			}}}},
			{Key: "category", Value: "$jobCategories"},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "companySize", Value: "$companySize"}, {Key: "category", Value: "$category"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetSkillsByExperienceLevel() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$mustSkills"}},
		{{Key: "$project", Value: bson.D{
			{Key: "experienceLevel", Value: bson.D{{Key: "$switch", Value: bson.D{
				{Key: "branches", Value: bson.A{
					bson.D{{Key: "case", Value: bson.D{{Key: "$lte", Value: bson.A{"$yearsOfExperience", 2}}}}, {Key: "then", Value: "Junior"}},
					bson.D{{Key: "case", Value: bson.D{{Key: "$lte", Value: bson.A{"$yearsOfExperience", 5}}}}, {Key: "then", Value: "Mid-Level"}},
				}},
				{Key: "default", Value: "Senior"},
			}}}},
			{Key: "skill", Value: "$mustSkills"},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "experienceLevel", Value: "$experienceLevel"}, {Key: "skill", Value: "$skill"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetCompaniesBySize() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$company"},
			{Key: "companySize", Value: bson.D{{Key: "$first", Value: "$companySize"}}},
			{Key: "jobCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "companySize", Value: -1}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "company", Value: "$_id"},
			{Key: "companySize", Value: 1},
			{Key: "jobCount", Value: 1},
			{Key: "_id", Value: 0},
		}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetCompaniesBySizeAndType(sizeType string) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sizeCriteria bson.D
	switch sizeType {
	case "small":
		sizeCriteria = bson.D{{Key: "$lte", Value: 50}}
	case "medium":
		sizeCriteria = bson.D{{Key: "$gt", Value: 50}, {Key: "$lte", Value: 250}}
	case "large":
		sizeCriteria = bson.D{{Key: "$gt", Value: 250}}
	default:
		return nil, errors.New("invalid size type")
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "companySize", Value: sizeCriteria}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "company", Value: "$company"}, {Key: "location", Value: "$location"}}},
			{Key: "companySize", Value: bson.D{{Key: "$first", Value: "$companySize"}}},
			{Key: "locationCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "locationCount", Value: -1}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id.company"},
			{Key: "companySize", Value: bson.D{{Key: "$first", Value: "$companySize"}}},
			{Key: "location", Value: bson.D{{Key: "$push", Value: "$_id.location"}}},
			{Key: "totalJobs", Value: bson.D{{Key: "$sum", Value: "$locationCount"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "companySize", Value: -1}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "company", Value: "$_id"},
			{Key: "companySize", Value: 1},
			{Key: "location", Value: 1},
			{Key: "totalJobs", Value: 1},
			{Key: "_id", Value: 0},
		}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetCompanySizeDistribution() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$facet", Value: bson.D{
			{Key: "sizes", Value: bson.A{
				bson.D{{Key: "$bucket", Value: bson.D{
					{Key: "groupBy", Value: "$companySize"},
					{Key: "boundaries", Value: bson.A{0, 51, 251, math.MaxInt64}},
					{Key: "default", Value: "Unknown"},
					{Key: "output", Value: bson.D{
						{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
					}},
				}}},
			}},
			{Key: "total", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "sizes", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: "$sizes"},
				{Key: "as", Value: "size"},
				{Key: "in", Value: bson.D{
					{Key: "sizeCategory", Value: bson.D{{Key: "$switch", Value: bson.D{
						{Key: "branches", Value: bson.A{
							bson.D{{Key: "case", Value: bson.D{{Key: "$eq", Value: bson.A{"$$size._id", "Unknown"}}}}, {Key: "then", Value: "Unknown"}},
							bson.D{{Key: "case", Value: bson.D{{Key: "$lte", Value: bson.A{"$$size._id", 50}}}}, {Key: "then", Value: "Small"}},
							bson.D{{Key: "case", Value: bson.D{{Key: "$lte", Value: bson.A{"$$size._id", 250}}}}, {Key: "then", Value: "Medium"}},
						}},
						{Key: "default", Value: "Large"},
					}}}},
					{Key: "count", Value: "$$size.count"},
					{Key: "percentage", Value: bson.D{{Key: "$divide", Value: bson.A{"$$size.count", bson.D{{Key: "$arrayElemAt", Value: bson.A{"$total.count", 0}}}}}}},
				}},
			}}}},
		}}},
		{{Key: "$unwind", Value: "$sizes"}},
		{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$sizes"}}}},
		{{Key: "$sort", Value: bson.D{{Key: "sizeCategory", Value: 1}}}},
	}

	results, err := s.aggregateResults(ctx, pipeline)
	if err != nil {
		return nil, apperrors.NewBaseError(
			apperrors.ErrCodeStorage,
			"Failed to execute aggregation query",
			err,
		)
	}

	return results, nil
}

func (s *JobStatisticsService) GetJobPostingsPerDay() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$dateToString", Value: bson.D{
				{Key: "format", Value: "%Y-%m-%d"},
				{Key: "date", Value: "$postingDate"},
			}}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "day", Value: "$_id"},
			{Key: "count", Value: 1},
			{Key: "_id", Value: 0},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "day", Value: 1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetJobPostingsPerMonth() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$dateToString", Value: bson.D{
				{Key: "format", Value: "%Y-%m"},
				{Key: "date", Value: "$postingDate"},
			}}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "month", Value: "$_id"},
			{Key: "count", Value: 1},
			{Key: "_id", Value: 0},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "month", Value: 1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetJobPostingsPerCompany() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$company"},
			{Key: "numberOfPostings", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "postingDates", Value: bson.D{{Key: "$push", Value: "$postingDate"}}},
			{Key: "postingUrls", Value: bson.D{{Key: "$push", Value: "$url"}}},
			{Key: "mostRecentPost", Value: bson.D{{Key: "$max", Value: "$postingDate"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "numberOfPostings", Value: -1}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "companyName", Value: "$_id"},
			{Key: "numberOfPostings", Value: 1},
			{Key: "postingDates", Value: 1},
			{Key: "postingUrls", Value: 1},
			{Key: "mostRecentPost", Value: 1},
			{Key: "_id", Value: 0},
		}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetMustSkillFrequencyPerDay(skill string) ([]bson.M, error) {
	return s.getSkillFrequencyPerDay(skill, "mustSkills")
}

func (s *JobStatisticsService) GetOptionalSkillFrequencyPerDay(skill string) ([]bson.M, error) {
	return s.getSkillFrequencyPerDay(skill, "optionalSkills")
}

func (s *JobStatisticsService) getSkillFrequencyPerDay(skill, skillField string) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := mongo.Pipeline{
		{{
			Key: "$match",
			Value: bson.D{{
				Key: skillField,
				Value: bson.D{{
					Key:   "$regex",
					Value: primitive.Regex{Pattern: "^" + skill + "$", Options: "i"},
				}},
			}},
		}},
		{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: bson.D{{
					Key: "$dateToString",
					Value: bson.D{
						{Key: "format", Value: "%Y-%m-%d"},
						{Key: "date", Value: "$postingDate"},
					},
				}}},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			},
		}},
		{{
			Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "day", Value: "$_id"},
				{Key: "count", Value: 1},
			},
		}},
		{{Key: "$sort", Value: bson.D{{Key: "day", Value: 1}}}},
	}
	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) GetJobCategoryCounts() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$jobCategories"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$jobCategories"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	return s.aggregateResults(ctx, pipeline)
}

func (s *JobStatisticsService) aggregateResults(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error) {
	return s.storage.AggregateJobs(ctx, pipeline)
}
